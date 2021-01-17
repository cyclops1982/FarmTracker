/* 
Simple utility to parse base64 string that we receive from Chripstack. 
The base64 string is the real data that we get from the lora device.
We can use the binary package to directly convert the []byte into a struct of the correct format.

This little program demonstrates that, and acts as utility to quickly see what the base64 content would really be.

*/
package main


import (
	"fmt"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"bytes"
	"github.com/cyclops1982/farmtracker/messagestructs"
)


func main() {
	var base = flag.String("base64", "l+j0X3YTbgKhHk7Awf8AAAAAIAYoDxIC", "some base64 stuff to decode")
	flag.Parse()

	fmt.Printf("base64 text: '%s'\n", *base)
	bs, err := base64.StdEncoding.DecodeString(*base)
	if err != nil {
		fmt.Printf("Failed to decode base64: %s\n", err)
	}
	fmt.Printf("Bytes: %q\n", bs)

	var lora loramsgs.SodaqUniversalTracker
	reader := bytes.NewReader(bs)
	err = binary.Read(reader, binary.LittleEndian, &lora)
	if err != nil {
		fmt.Println("Failed to read binary data: ", err)
	}
	fmt.Printf("Unixtime: %d\nVoltage: %d\nLat/Long: %q/%q\n", lora.Unixtime, lora.RawVoltage, lora.Latitude, lora.Longitude)
}
