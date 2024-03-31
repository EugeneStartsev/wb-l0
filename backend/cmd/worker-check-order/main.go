package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"wb/backend/cache"
	"wb/backend/nats/publisher"
	subscriber "wb/backend/nats/subscriber"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	dbConn := flag.String("db",
		"host=localhost dbname=db user=wb password=admin sslmode=disable",
		"database connection string")
	httpPort := flag.Int("http-port", 4000, "HTTP API port")
	flag.Parse()

	postgres, err := sql.Open("postgres", *dbConn)
	if err != nil {
		log.Fatal(err)
	}

	defer func(postgres *sql.DB) {
		err = postgres.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(postgres)

	db := goqu.New("postgres", postgres)
	lru := cache.New(100)

	err = cache.RecoverLruFromPostgres(db, lru)
	if err != nil {
		log.Fatal(err)
	}

	pub := publisher.StartPublisher()

	sub := subscriber.StartSubscriber(lru)

	s := newHttpServer(db, lru)

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
