package publisher

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit"
	"github.com/nats-io/stan.go"
	"log"
	"time"
	"wb/backend/structs"
)

func StartPublisher() stan.Conn {
	var ord structs.Orders

	sc, err := stan.Connect("test-cluster", "publisher")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 3)

			gofakeit.Struct(&ord)

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
