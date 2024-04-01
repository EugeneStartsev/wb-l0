package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"wb/backend/cache"
	"wb/backend/postgres"
	"wb/backend/structs"
)

type httpServer struct {
	storage postgres.Storage
	router  *gin.Engine
	lru     *cache.LRU
}

func newHttpServer(storage postgres.Storage, lru *cache.LRU) *httpServer {
	s := httpServer{
		storage: storage,
		router:  gin.Default(),
		lru:     lru,
	}

	s.router.GET("/order", s.handleGetOrder)
	s.router.POST("/order/create", s.handlePostOrder)

	return &s
}

func (s *httpServer) run(listenAddr string) error {
	return s.router.Run(listenAddr)
}

func (s *httpServer) handleGetOrder(c *gin.Context) {
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

	orders, err := s.storage.GetOrder(query.Uid)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("не удалось получить заказ"))
		return
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

	marshalData, err := json.Marshal(data)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = s.storage.SaveOrder(data)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	s.lru.Set(data.ID, marshalData)
}
