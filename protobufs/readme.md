# Protobufs

This folder contains protobuf definitions used in the FarmTracker project.

Compile the protobuf with ```
protoc --go_out=. --go_opt=module=github.com/cyclops1982/farmtracker/protobufs protobufs.proto
```. This will make the file in the same directory.

## Protoc & protoc-gen-go

Protoc wasn't super easy to get going. No clue why this is not aligned with the golang installation manuel. What we did:
```
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip
mkdir -p protoc && cd protoc
unzip ../protoc-3.14.0-linux-x86_64.zip
chmod 755 * -R
cp bin/protoc /usr/local/go/bin/
cp -r bin/include /usr/local/go/
```

Then, to install protoc-gen-go:
```
go install google.golang.org/protobuf/cmd/protoc-gen-go
```
and add it to the PATH
