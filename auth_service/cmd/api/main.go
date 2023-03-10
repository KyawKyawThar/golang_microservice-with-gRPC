package main

import (
	"auth_service/cmd/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var counts int64

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("---------------------------------------------")
	log.Println("Attempting to connect to Postgres...")

	conn := OpenDB()

	if conn == nil {
		log.Panic("can't connect to postgres!")
	}

	app := Config{DB: conn, Models: data.New(conn)}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	log.Printf("Starting authentication end service on port %s\n", webPort)

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func OpenDB() *sql.DB {

	// connect to postgres
	dsn := os.Getenv("DSN")

	for {
		connection, err := ConnectDB(dsn)

		if err != nil {
			log.Println("Postgres not ready...")
			counts++
		} else {
			log.Println("Connected to database!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue

	}
}

func ConnectDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
