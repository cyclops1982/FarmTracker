package main

import (
	"fmt"
	"log"
	"flag"
	"context"
	"database/sql"
	"time"
	_ "github.com/go-sql-driver/mysql"
)


var pool *sql.DB


func CheckConnection(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	err := pool.Ping();
	if err != nil {
		log.Fatalf("Ping() failed - unable to connect to database: %v", err)
	}
}

func main() {
	var dsn  = flag.String("sqlconstring", "farmtracker:MyGreatPassword@tcp(localhost)/FarmTracker", "The DSN Connection String to use to connect to the MySQL DB.")
	flag.Parse()

//	lat := 22.23232
//	long := 55.12123
//	devId := 1
	var err error
	pool, err = sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatal("Failed to SQL: ", err)
	}
	defer pool.Close()


	pool.SetConnMaxLifetime(time.Minute * 3)
	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(10)

	// Create context
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	CheckConnection(ctx)
	err = pool.Ping()

	insert, err := pool.Prepare("INSERT INTO Location(DeviceId, Location) VALUES((SELECT Id FROM Device WHERE DeviceEUI = ?), ST_GeomFromText(?))")
	defer insert.Close()
	if err != nil {
		log.Fatal("Failed to prepare: ", err)
	}

	devEUI := "0004a30b00eb5e28"
	long := 12.122223
	lat := 40.23233
	rows, err := insert.Exec(devEUI, fmt.Sprintf("POINT(%f %f)", long, lat))
	if err != nil {
		log.Fatal("Failed to insert: ", err)
	}
	log.Println("Added ", rows, "rows.")
}
