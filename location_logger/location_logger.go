package main

import (
	"log"
	"net"
	"encoding/json"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/cyclops1982/farmtracker/messagestructs"
	"github.com/cyclops1982/farmtracker/enhancedconn"
	"bytes"
	"time"
	"context"
	"database/sql"
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
	//defer sqlInsertLocation.Close()
	//defer sqlInsertBattery.Close()

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

	// Prepare the insert statement
	sqlInsert, err := pool.Prepare("INSERT INTO Location(LoggedOn, DeviceId, Location) VALUES(?, (SELECT Id FROM Device WHERE DeviceEUI = ?), ST_GeomFromText(?))")
	defer sqlInsert.Close()
	if err != nil {
		log.Fatalf("Failed to prepare INSERT statement: %v\n", err)
	}

	unixtimebytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(unixtimebytes, uint64(*fromUnixtime))
	con.Write(unixtimebytes)
	con.Write([]byte("up."))
	msgLength := make([]byte, 2)
	for {
		nBytes, err := con.Read(msgLength)
		msgLengthUint16 := binary.BigEndian.Uint16(msgLength)
		log.Printf("Length of message that's coming: %d\n", msgLengthUint16)
		if nBytes != 2 {
			log.Fatal("We really expect 2 bytes for a messagelength.")
		}
		msgData := make([]byte, int(msgLengthUint16))
		
		nBytes, err = con.ReadBytes(msgData)
		if err != nil {
			log.Println("Failed to read:", err)
			continue
		}
		
		// Convert received stuff to JSON.
		//TODO: make this something more strongly typed - need to check what happens if our struct is not 100% aligned.
		var jsonData interface{}
		err = json.Unmarshal(msgData, &jsonData)
		if err != nil {
			log.Println("Failed to parse JSON:", err)
			continue
		}
		// get the properties that we'd like to have.
		realData := jsonData.(map[string]interface{})
		devEUI, ok := realData["devEUI"].(string)
		if ok == false {
			log.Println("Failed to convert data to string. Skipping.")
			continue
		}
		log.Println("DEVEUI: ",devEUI)
		base64data, ok := realData["data"].(string)
		if ok == false {
			log.Println("Failed to convert data to string. Skipping.")
			continue
		}

		// convert the base64 string to a []byte
		bs, err := base64.StdEncoding.DecodeString(base64data)
		if err != nil {
			log.Printf("Failed to get decode base64 string '%s'. Skipping.\n", base64data)
			continue
		}
		var loraMsg loramsgs.SodaqUniversalTracker
		byteReader := bytes.NewReader(bs)
		err = binary.Read(byteReader, binary.LittleEndian, &loraMsg)
		if err != nil {
			log.Printf("Couldn't unpack binary array from base64 data ('%s') into Lora Msg Struct. Skipping.\n", base64data)
			continue
		}

		log.Println("Battery Status: ", loraMsg.RawVoltage)
		go AddRecords(pool, devEUI, &loraMsg)

	}
}
