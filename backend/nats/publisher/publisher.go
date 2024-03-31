package publisher

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/signal"
	"time"
	"wb/backend/structs"
)

func StartPublisher() {
	var ord structs.Orders

	sc, err := stan.Connect("test-cluster", "publisher")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 3)

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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Wait for the signal.
	<-sigCh

	sc.Close()
}
