package nats

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"wb/backend/cache"
	"wb/backend/structs"
)

func startSubscriber(lru *cache.LRU) {
	sc, err := stan.Connect("test-cluster", "subscriber")
	if err != nil {
		log.Fatalf("Subscriber: %s", err)
	}

	data := *new(structs.Orders)

	sub, err := sc.Subscribe("JsonPipe", func(m *stan.Msg) {
		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			log.Println(err)
		} else {
			// добавление в кэш
			lru.Set(data.ID, m.Data)
			// добавление в бд
			err = storage.SaveOrder(data)
			if err != nil {
				log.Println(err)
			}
		}
	}, stan.DeliverAllAvailable())
}
