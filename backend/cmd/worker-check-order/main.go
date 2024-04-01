package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"wb/backend/cache"
	"wb/backend/nats/publisher"
	subscriber "wb/backend/nats/subscriber"
	"wb/backend/postgres"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	dbConn := flag.String("db",
		"host=localhost dbname=db user=wb password=admin sslmode=disable",
		"database connection string")
	httpPort := flag.Int("http-port", 4000, "HTTP API port")
	flag.Parse()

	storage, err := postgres.NewDB(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	lru := cache.New(100)

	err = cache.RecoverLruFromPostgres(storage, lru)
	if err != nil {
		log.Fatal(err)
	}

	pub := publisher.StartPublisher()

	sub := subscriber.StartSubscriber(lru, storage)

	s := newHttpServer(storage, lru)

	go func() {
		err = s.run(fmt.Sprintf(":%d", *httpPort))
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	<-sigch

	pub.Close()
	sub.Unsubscribe()
	sub.Close()
}
