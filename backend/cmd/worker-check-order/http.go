package main

import (
	"encoding/json"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"wb/backend/structs"
)

type httpServer struct {
	db     *goqu.Database
	router *gin.Engine
	lru    *LRU
}

func newHttpServer(db *goqu.Database) *httpServer {
	s := httpServer{
		db:     db,
		router: gin.Default(),
		lru:    New(100),
	}

	err := s.recoverLruFromPostgres()
	if err != nil {
		log.Fatal(err)
	}

	s.router.GET("/order", s.handleGetOrder)
	s.router.POST("/order/create", s.handlePostOrder)

	return &s
}

func (s *httpServer) run(listenAddr string) error {
	return s.router.Run(listenAddr)
}

func (s *httpServer) handleGetOrder(c *gin.Context) {
	var data structs.Ord
	var items []structs.Item

	var query struct {
		Uid string `form:"uid"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if query.Uid == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if val, ok := s.lru.Get(query.Uid); ok {
		var ord structs.Orders

		err := json.Unmarshal(val, &ord)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, ord)
		return
	}

	found, err := s.db.From("orders").
		InnerJoin(goqu.T("delivery"), goqu.Using("uid")).
		InnerJoin(goqu.T("payment"), goqu.Using("uid")).
		Where(goqu.Ex{"uid": query.Uid}).
		ScanStruct(&data)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !found {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = s.db.From("items").Where(goqu.Ex{"uid": query.Uid}).ScanStructs(&items)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var orders = structs.Orders{
		ID:                data.ID,
		TrackNumber:       data.TrackNumber,
		Entry:             data.Entry,
		Delivery:          data.Delivery,
		Payments:          data.Payment,
		Items:             items,
		Locale:            data.Locale,
		InternalSignature: data.InternalSignature,
		CustomerID:        data.CustomerID,
		DeliveryService:   data.DeliveryService,
		ShardKey:          data.ShardKey,
		SmID:              data.SmID,
		DateCreated:       data.DateCreated,
		OofShard:          data.OofShard,
	}

	c.JSON(http.StatusOK, orders)
}

func (s *httpServer) handlePostOrder(c *gin.Context) {
	var data structs.Orders

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tx, err := s.db.Begin()

	_, err = tx.Insert("orders").Rows(goqu.Record{
		"uid":                data.ID,
		"track_number":       data.TrackNumber,
		"entry":              data.Entry,
		"locale":             data.Locale,
		"internal_signature": data.InternalSignature,
		"customer_id":        data.CustomerID,
		"delivery_service":   data.DeliveryService,
		"shardkey":           data.ShardKey,
		"sm_id":              data.SmID,
		"date_created":       data.DateCreated,
		"oof_shard":          data.OofShard,
	}).Executor().Exec()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	_, err = tx.Insert("delivery").Rows(goqu.Record{
		"uid":     data.ID,
		"phone":   data.Delivery.PhoneNumber,
		"zip":     data.Delivery.Zip,
		"city":    data.Delivery.City,
		"address": data.Delivery.Address,
		"region":  data.Delivery.Region,
		"email":   data.Delivery.Email,
		"name":    data.Delivery.Name,
	}).Executor().Exec()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	_, err = tx.Insert("payment").Rows(goqu.Record{
		"uid":           data.ID,
		"transaction":   data.Payments.Transaction,
		"request_id":    data.Payments.RequestID,
		"currency":      data.Payments.Currency,
		"provider":      data.Payments.Provider,
		"amount":        data.Payments.Amount,
		"payment_dt":    data.Payments.PaymentDT,
		"bank":          data.Payments.BankName,
		"delivery_cost": data.Payments.Cost,
		"goods_total":   data.Payments.TotalGoods,
		"custom_fee":    data.Payments.CustomFee,
	}).Executor().Exec()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	for _, item := range data.Items {
		_, err = tx.Insert("items").Rows(goqu.Record{
			"uid":          data.ID,
			"chrt_id":      item.ChrtID,
			"track_number": item.TrackNumber,
			"price":        item.Price,
			"rid":          item.RID,
			"name":         item.Name,
			"sale":         item.Sale,
			"size":         item.Size,
			"total_price":  item.TotalPrice,
			"nm_id":        item.NmID,
			"brand":        item.Brand,
			"status":       item.Status,
		}).Executor().Exec()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	marshalData, err := json.Marshal(data)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.lru.Set(data.ID, marshalData)
}

func getOrdersFromPostgres(db *goqu.Database) []structs.Orders {
	var ords []structs.Ord

	err := db.From("orders").
		InnerJoin(goqu.T("delivery"), goqu.Using("uid")).
		InnerJoin(goqu.T("payment"), goqu.Using("uid")).
		ScanStructs(&ords)
	if err != nil {
		return nil
	}

	res := make([]structs.Orders, 0, len(ords))

	for _, val := range ords {
		var items []structs.Item
		err = db.From("items").Where(goqu.Ex{"uid": val.ID}).ScanStructs(&items)
		if err != nil {
			return nil
		}

		var orders = structs.Orders{
			ID:                val.ID,
			TrackNumber:       val.TrackNumber,
			Entry:             val.Entry,
			Delivery:          val.Delivery,
			Payments:          val.Payment,
			Items:             items,
			Locale:            val.Locale,
			InternalSignature: val.InternalSignature,
			CustomerID:        val.CustomerID,
			DeliveryService:   val.DeliveryService,
			ShardKey:          val.ShardKey,
			SmID:              val.SmID,
			DateCreated:       val.DateCreated,
			OofShard:          val.OofShard,
		}

		res = append(res, orders)
	}

	return res
}

func (s *httpServer) recoverLruFromPostgres() error {
	orders := getOrdersFromPostgres(s.db)
	if orders == nil {
		return fmt.Errorf("не удалось восстановить кэш и бд")
	}

	for _, val := range orders {
		marshalOrder, err := json.Marshal(val)
		if err != nil {
			return err
		}

		s.lru.Set(val.ID, marshalOrder)
	}

	return nil
}
