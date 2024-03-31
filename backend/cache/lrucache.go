package cache

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"sync"
	"wb/backend/structs"
)

type order struct {
	key   string
	order []byte
}

type LRU struct {
	rw       sync.RWMutex
	capacity int
	items    map[string]*list.Element
	queue    *list.List
}

func New(cap int) *LRU {
	return &LRU{
		capacity: cap,
		items:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

func (l *LRU) Get(key string) ([]byte, bool) {
	elem, isExist := l.items[key]

	if isExist == false {
		return nil, false
	}
	l.rw.RLock()
	defer l.rw.RUnlock()

	l.queue.MoveToFront(elem)
	return elem.Value.(*order).order, true
}

func (l *LRU) Set(key string, value []byte) {
	l.rw.RLock()
	defer l.rw.Unlock()

	if elem, isExist := l.items[key]; isExist == true {
		l.queue.MoveToFront(elem)
		elem.Value.(*order).order = value
	}

	if l.queue.Len() == l.capacity {
		l.purge()
	}

	item := &order{
		key:   key,
		order: value,
	}

	elem := l.queue.PushFront(item)
	l.items[key] = elem
}

func (l *LRU) Clear() {
	l.queue = list.New()
}

func RecoverLruFromPostgres(db *goqu.Database, lru *LRU) error {
	orders := getOrdersFromPostgres(db)
	if orders == nil {
		return fmt.Errorf("не удалось восстановить кэш и бд")
	}

	for _, val := range orders {
		marshalOrder, err := json.Marshal(val)
		if err != nil {
			return err
		}

		lru.Set(val.ID, marshalOrder)
	}

	return nil
}

func (l *LRU) purge() {
	if elem := l.queue.Back(); elem != nil {
		deletedOrder := l.queue.Remove(elem).(*order)
		delete(l.items, deletedOrder.key)
	}
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
