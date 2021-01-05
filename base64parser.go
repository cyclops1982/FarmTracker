package main


import (
	"fmt"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"bytes"
)

type LoraMsg struct {
	Unixtime uint32
	RawVolt uint8
	BoardTemp int8
	Lat int32
	Long int32
}

func main() {
	var base = flag.String("base64", "l+j0X3YTbgKhHk7Awf8AAAAAIAYoDxIC", "some base64 stuff to decode")
	flag.Parse()

	fmt.Printf("base64 text: '%s'\n", *base)
	bs, err := base64.StdEncoding.DecodeString(*base)
	if err != nil {
		fmt.Printf("Failed to decode base64: %s\n", err)
	}
	fmt.Printf("Bytes: %q\n", bs)
	i := binary.LittleEndian.Uint32(bs[0:4])
	fmt.Printf("In int: %d\n", i)

	var lora LoraMsg
	reader := bytes.NewReader(bs)
	err = binary.Read(reader, binary.LittleEndian, &lora)
	if err != nil {
		fmt.Println("Failed to read binary data: ", err)
	}
	fmt.Printf("Unixtime: %d\nVoltage: %d\nLat/Long: %q/%q\n", lora.Unixtime, lora.RawVolt, lora.Lat, lora.Long)
}
