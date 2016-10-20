/*
This example demonstrates:

1. Change Lock Pin
2. Unlock of kinetic device.

Default pin is blank string, for PIN operation, SSL connection is required.
*/
package main

import (
	"fmt"
	"os"

	kinetic "github.com/yongzhy/kinetic-go"
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

	status, err := conn.SetLockPin([]byte(""), []byte("PIN"))
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Blocking SetLockPin Failure: ", err, status.String())
		os.Exit(-1)
	}

	status, err = conn.UnlockDevice([]byte("PIN"))
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Blocking UnlockDevice Failure: ", err, status.String())
		os.Exit(-1)
	}

}
