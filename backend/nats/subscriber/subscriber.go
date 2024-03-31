package nats

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"wb/backend/cache"
	"wb/backend/structs"
)

func StartSubscriber(lru *cache.LRU) stan.Subscription {
	sc, err := stan.Connect("test-cluster", "subscriber")
	if err != nil {
		log.Fatalf("Subscriber: %s", err)
	}

	data := *new(structs.Orders)

	sub, err := sc.Subscribe("Json", func(m *stan.Msg) {
		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			log.Println(err)
		} else {
			lru.Set(data.ID, m.Data)
		}
	}, stan.DeliverAllAvailable())

	return sub
}
