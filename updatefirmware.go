/*
This example demonstrates:

Update kinetic drive firmware.

*/

package main

import (
	"fmt"
	"flag"

	kinetic "github.com/Kinetic/kinetic-go"
	"os"
)

var host = flag.String("host", "127.0.0.1", "Kinetic device IP address")
var file = flag.String("lod", "", "New firmware file path")

func main() {
	flag.Parse()

	if file == nil || *file == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Set the log leverl to debug
	kinetic.SetLogLevel(kinetic.LogLevelDebug)

	// Client options
	var option = kinetic.ClientOptions{
		Host: *host, // Need to test with drive
		Port: 8123,
		User: 1,
		Hmac: []byte("asdfasdf")}

	conn, err := kinetic.NewBlockConnection(option)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = kinetic.UpdateFirmware(conn, *file)
	if err != nil {
		fmt.Println("Firmware update fail: ", *file, err)
	}
}
