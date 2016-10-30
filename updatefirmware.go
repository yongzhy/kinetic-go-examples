/*
This example demonstrates:

Update kinetic drive firmware.

*/

package main

import (
	"fmt"

	kinetic "github.com/Kinetic/kinetic-go"
)

func main() {
	// Set the log leverl to debug
	kinetic.SetLogLevel(kinetic.LogLevelDebug)

	// Client options
	var option = kinetic.ClientOptions{
		Host: "127.0.0.1", // Need to test with drive
		Port: 8123,
		User: 1,
		Hmac: []byte("asdfasdf")}

	conn, err := kinetic.NewBlockConnection(option)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	file := "/change/this/to/correct/path/AD-installer-v44.01.04.slod"
	err = kinetic.UpdateFirmware(conn, file)
	if err != nil {
		fmt.Println("Firmware update fail: ", file, err)
	}
}
