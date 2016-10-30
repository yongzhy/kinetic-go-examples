/*
This example demonstrates:

1. Change Erase Pin
2. InstantErase of kinetic device.

Default pin is blank string, for PIN operation, SSL connection is required.
*/

package main

import (
	"fmt"
	"os"

	kinetic "github.com/Kinetic/kinetic-go"
)

var option = kinetic.ClientOptions{
	Host: "127.0.0.1", // Test with Simulator
	//Port: 8123,
	Port:   8443, // For SSL connection
	User:   1,
	Hmac:   []byte("asdfasdf"),
	UseSSL: true,
}

func main() {
	kinetic.SetLogLevel(kinetic.LogLevelDebug)
	conn, err := kinetic.NewBlockConnection(option)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	status, err := conn.SetErasePin([]byte(""), []byte("PIN"))
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Blocking SetErasePin Failure: ", err, status.String())
		os.Exit(-1)
	}

	status, err = conn.InstantErase([]byte("PIN"))
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Blocking InstantErase Failure: ", err, status.String())
		os.Exit(-1)
	}

	status, err = conn.SecureErase([]byte(""))
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Blocking SecureErase Failure: ", err, status.String())
		os.Exit(-1)
	}

}
