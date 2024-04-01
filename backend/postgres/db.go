package postgres

import (
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"wb/backend/structs"
)

type Storage struct {
	db *goqu.Database
}

func NewDB(dbConn *string) (Storage, error) {
	postgres, err := sql.Open("postgres", *dbConn)
	if err != nil {
		return Storage{}, err
	}

	return Storage{goqu.New("postgres", postgres)}, nil
}

func (s *Storage) GetOrder(uid string) (structs.Orders, error) {
	var data structs.Ord
	var items []structs.Item

	found, err := s.db.From("orders").
		InnerJoin(goqu.T("delivery"), goqu.Using("uid")).
		InnerJoin(goqu.T("payment"), goqu.Using("uid")).
		Where(goqu.Ex{"uid": uid}).
		ScanStruct(&data)
	if err != nil {
		return structs.Orders{}, err
	}

	if !found {
		return structs.Orders{}, err
	}

	err = s.db.From("items").Where(goqu.Ex{"uid": uid}).ScanStructs(&items)
	if err != nil {
		return structs.Orders{}, err
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

	return orders, nil
}

func (s *Storage) SaveOrder(data structs.Orders) error {
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
		tx.Rollback()
		return err
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
		tx.Rollback()
		return err
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
		tx.Rollback()
		return err
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
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetOrdersFromPostgres() []structs.Orders {
	var ords []structs.Ord

	err := s.db.From("orders").
		InnerJoin(goqu.T("delivery"), goqu.Using("uid")).
		InnerJoin(goqu.T("payment"), goqu.Using("uid")).
		ScanStructs(&ords)
	if err != nil {
		return nil
	}

	res := make([]structs.Orders, 0, len(ords))

	for _, val := range ords {
		var items []structs.Item
		err = s.db.From("items").Where(goqu.Ex{"uid": val.ID}).ScanStructs(&items)
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
