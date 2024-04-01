package publisher

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/nats-io/stan.go"
	"log"
	"time"
	"wb/backend/structs"
)

func StartPublisher() stan.Conn {
	var ord structs.Orders

	sc, err := stan.Connect("test-cluster", "publisher")
	if err != nil {
		log.Fatalf("Publisher: %s", err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 10)

			err = gofakeit.Struct(&ord)
			if err != nil {
				log.Println(err)
			}

			jsonToSend, err := json.Marshal(ord)
			if err != nil {
				log.Println(err)
			}

			err = sc.Publish("Json", jsonToSend)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	return sc
}
