# Farmtracker

FarmTracker is a project that helps farmers gain insight to their farm and assets on their farm.
The initial goals are:

-  Sheep/Animal tracking
-  Water sensors

## Technical notes

Go Version:

```
$ go version
go version go1.15.6 linux/amd64
```

## Architecture

The following tools/apps exist:

-  http_receiver.go - This is a web service that receives a POST message from Chirpstack's http integration. Does validation if it is what we expect (a POST to a specific URL). Messages are written to disk in a specific structure that dir_processor.go reads.
-  dir_processor.go - Simple application that waits for something to connect, expects what to send out and then sends it. "Monitors" the filesystem for new files that are added. If file matches the filter, will then send it over.
-  base64parser.go - Simple utility to output the base64 data in readable format. Also a test to see if our base64 actually unpacks correctly.
-  location_logger.go - receives the 'up' messages from the dir_processor over the network and pushes it into the db.

## Installation/Running

-  Run the db_create.sql script. Please check the content.
-  Run the table_create.sql script.
-  Start the http_receiver.
-  Start the dir_processor.

## TODO's

-  Make the website a PWA by adding a manifest file
-  Add 'magic bytes' to the data communication between dir_processor and it's clients.
-  Enhance the location_logger to log out the last unixtime that it received.
-  Enhance the location_logger to be more failsafe on what it reads
-  Add dir_processor client to write the received JSON to a DB, so that we can query that.
-  Add dir_processor client that writes to duckDB for time series info
-  Add graph of time-series data to web app

### Query notes

-  MSG per day:
   SELECT COUNT(LoggedOn), DATE_FORMAT(LoggedOn, "%Y-%m-%d") AS YMD from Location GROUP BY YMD;

## Other notes

-  Great way to demonstrate leaflet layers: https://leaflet-extras.github.io/leaflet-providers/preview/
-  Nice table implementation: https://codepen.io/geoffyuen/pen/FCBEg
-  jQuery is not needed in most cases: http://youmightnotneedjquery.com/ - let's try to keep it lightweight.
-  Emojis https://github.com/twitter/twemoji (we could use these as images on the map?)
