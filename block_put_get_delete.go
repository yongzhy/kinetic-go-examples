/*
This example demonstrates:

1. Make blocking connection to Kinetic drive
2. PUT an object
3. Get back the object
4. Delete the object

*/

package main

import (
	"bytes"
	"fmt"

	kinetic "github.com/Kinetic/kinetic-go"
)

var option = kinetic.ClientOptions{
	Host: "127.0.0.1", // Test with Simulator
	Port: 8123,
	User: 1,
	Hmac: []byte("asdfasdf")}

func main() {
	kinetic.SetLogLevel(kinetic.LogLevelDebug)

	conn, err := kinetic.NewBlockConnection(option)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Put
	pentry := kinetic.Record{
		Key:   []byte("Test Object"),
		Value: []byte("Test Object Data"),
		Sync:  kinetic.SYNC_WRITETHROUGH,
		Algo:  kinetic.ALGO_SHA1,
		Tag:   []byte(""),
		Force: true,
	}
	status, err := conn.Put(&pentry)
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Blocking Put Failure", err, status)
	}

	gentry, status, err := conn.Get(pentry.Key)
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Blocking Get Failure", err, status)
	}

	if !bytes.Equal(pentry.Key, gentry.Key) {
		fmt.Printf("Key Mismatch: [%s] vs [%s]\n", pentry.Key, gentry.Key)
	}
	if !bytes.Equal(pentry.Value, gentry.Value) {
		fmt.Printf("Value Mismatch: [%s] vs [%s]\n", pentry.Value, gentry.Value)
	}

	dentry := kinetic.Record{
		Key:   pentry.Key,
		Sync:  pentry.Sync,
		Force: true,
	}
	status, err = conn.Delete(&dentry)
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Blocking Delete Failure", err, status)
	}

}
