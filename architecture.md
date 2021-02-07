# Architecture

This document aims to describe the architecture of the FarmTracker solution.
It should help with understanding how the projects work, and also set out some details of the internal protocols used.

## High level

The solution has a number of small applications that are listed here:

- `http_receiver` - This web service receives POST messages that are coming from Chirpstack's HTTP integration. It validates some details, but not a lot. The primary goals is to write the content of the POST request into a file on disk.
- `dir_processor` - This service scans the directories for new files, and if a new file is found will send it to it's connected clients. Clients that connect need to indicate when they were last connected and will receive all files since that date/time.
- `client_mysql_logger` - Retrieves data from the dir_processor and stores it in a DB. Performs the bulk of the processing work.
- `http_frontend` - the web application for the solution. Hosts all the (mostly) static files and has a REST api for the JavaScript to read data.

This architecture is partially based on https://apenwarr.ca/log/20190216

## Protocols

Between the various elements, there are some assumptions and expectations. These are written down here.
The aim is to keep a simple protocol that doesn't require loads of processing.

### dir_processor

`dir_processor` listens on a TCP port. Once a client is connected, it expects some data to know what it should send. The client indicates what data they would want to receive.

#### Initial Receiving data.
Once a client is connected, 8 bytes are read into an array and (little-endian) converted to `uint64`. This represents a unixtime. The reason for 64 bits is primarily because that's how golang represents it. The unixtime is used to send the messages *from* that timestamp.

After the initial 8 bytes for the unixtime, another 30 bytes are read. These 30 bytes should contain a dot ('.') character. This is used as text filter on the to-be received data. This effectively makes the client subscribe to a specific topic. If only a dot ('.') is send, then all data is send.

After these initial reads by `dir_processor`, it will start sending data.

#### Message format

- Investigating protobuf :)
