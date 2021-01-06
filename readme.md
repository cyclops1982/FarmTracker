# Farmtracker

FarmTracker is a project that helps farmers gain insight to their farm and assets on their farm.
The initial goals are:
- Sheep/Animal tracking
- Water sensors


## Technical notes

Go Version:
```
$ go version
go version go1.15.6 linux/amd64
```

## Architecture

The following tools/apps exist:
- http_receiver.go - This is a web service that receives a POST message from Chirpstack's http integration. Does validation if it is what we expect (a POST to a specific URL). Messages are written to disk in a specific structure that dir_processor.go reads.
- dir_processor.go - Simple application that waits for something to connect, expects what to send out and then sends it. "Monitors" the filesystem for new files that are added. If file matches the filter, will then send it over.
- base64parser.go - Simple utility to output the base64 data in readable format. Also a test to see if our base64 actually unpacks correctly.
- location_logger.go - receives the 'up' messages from the dir_processor over the network and pushes it into the db.

## Installation/Running

- Run the db_create.sql script. Please check the content.
- Run the table_create.sql script.
- Start the http_receiver.
- Start the dir_processor.

