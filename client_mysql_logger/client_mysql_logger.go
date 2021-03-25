package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cyclops1982/farmtracker/enhancedconn"
	loramsgs "github.com/cyclops1982/farmtracker/loramsgstructs"
	"github.com/cyclops1982/farmtracker/protobufs"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	_ "github.com/go-sql-driver/mysql"
)


func CreateConnection(dsn *string) *sql.DB {
	var pool *sql.DB
	var err error
	// Setup SQL connection
	pool, err = sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("Failed to connect to SQL: %v\n", err)
	}
	pool.SetConnMaxLifetime(time.Minute * 3)
	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(10)

	return pool
}

func CheckConnection(ctx context.Context, pool *sql.DB) {
	ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		log.Fatalf("PingContext() failed - Unable to connect to database: %v\n", err)
	}

	if err := pool.Ping(); err != nil {
		log.Fatalf("Ping() failed - unable to connect to databsae: %v\n", err)
	}
}

func AddRecords(pool *sql.DB, deveui string, loraMsg *loramsgs.SodaqUniversalTracker) {
	var err error
	sqlInsertLocation, _ := pool.Prepare("INSERT INTO Location(LoggedOn, DeviceId, Location) VALUES(?, ?, ST_GeomFromText(?))")
	sqlInsertBattery, _ := pool.Prepare("INSERT INTO BatteryStatus(LoggedOn, DeviceId, RawValue) VALUES(?, ?, ?)")
	defer sqlInsertLocation.Close()
	defer sqlInsertBattery.Close()

	var deviceId int
	err = pool.QueryRow("SELECT Id FROM Device WHERE DeviceEUI=?", deveui).Scan(&deviceId)
	if err != nil {
		log.Printf("Failed to retrieve Id for Device '%s', error was: %v\n", deveui, err)
		return
	}

	long := float32(loraMsg.Longitude)/10000000
	lat := float32(loraMsg.Latitude)/10000000
	msgTime := time.Unix(int64(loraMsg.Unixtime), 0)

	_, err = sqlInsertLocation.Exec(msgTime, deviceId, fmt.Sprintf("POINT(%f %f)", lat, long))
	if err != nil {
		log.Printf("FAILED to insert location into DB: %v\n",err)
	}
	_, err = sqlInsertBattery.Exec(msgTime, deviceId, loraMsg.RawVoltage)
	if err != nil {
		log.Printf("FAILED to insert batterystatus into DB: %v\n",err)
	}
}

func main() {
	var err error
	var ipAddress = flag.String("server", "127.0.0.1", "The IP address (or hostname) of the server to connect to.")
	var tcpPort = flag.Int("port", 29000, "The port to use for the server connection.")
	var sqlConString  = flag.String("sqlconstring", "farmtracker:MyGreatPassword@tcp(localhost)/FarmTracker", "The DSN Connection String to use to connect to the MySQL DB.")
	//TODO: Add a 'from' parameter that just takes an amount of hours
	var fromUnixtime = flag.Int64("fromUnixtime", 0, "Set the unix timestamp from which we should receive messages.")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *ipAddress, *tcpPort)
	
	tmpcon, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to connect to %s. Exiting.", addr)
	}
	con := enhancedconn.EnhancedConn{tmpcon}


	// Create SQL connection & context
	pool := CreateConnection(sqlConString)
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	CheckConnection(ctx, pool)

	// tell the server what we want to receive.
	req := &protobufs.MessagesRequest{}
	req.DataToGet = protobufs.MessagesRequest_LoraUpdatesV1
	req.DataSince = timestamppb.New(time.Unix(*fromUnixtime, 0))

	err = con.SendProtobufMsg(req);
	if err != nil {
		log.Fatalf("Couldn't send protobuf to indicate what we want to receive: %v\n", err)
	}	
	
	var nbytes int
	var msgLength uint16
	for {
		// Read how long our message will be.
		msgLength, err = con.ReadLength()
		if err !=nil {
			log.Printf("Failed to read message size. Skipping and waiting for next bytes. Error: %v\n", err)
			continue
		}

		// read the actual message
		dataBytes := make([]byte, msgLength)
		nbytes, err = con.ReadBytes(dataBytes, 0)
		if err != nil {
			log.Printf("Read %d of %d expected bytes. Skipping to next message. Error: %v.\n", nbytes, msgLength, err)
			continue
		}

		
	}
}
