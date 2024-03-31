package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	"log"
	"wb/backend/cache"
	"wb/backend/nats/publisher"
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

	publisher.StartPublisher()

	s := newHttpServer(db, lru)

	err = s.run(fmt.Sprintf(":%d", *httpPort))
	if err != nil {
		log.Fatal(err)
	}
}
