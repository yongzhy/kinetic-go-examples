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

	kinetic "github.com/yongzhy/kinetic-go"
)

var option = kinetic.ClientOptions{
	Host: "127.0.0.1", // Test with Simulator
	Port: 8123,
	User: 1,
	Hmac: []byte("asdfasdf")}

func main() {

	// Set the log leverl to debug
	kinetic.SetLogLevel(kinetic.LogLevelDebug)

	conn, err := kinetic.NewNonBlockConnection(option)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// PUT
	pentry := kinetic.Record{
		Key:   []byte("Test Object"),
		Value: []byte("Test Object Data"),
		Sync:  kinetic.SYNC_WRITETHROUGH,
		Algo:  kinetic.ALGO_SHA1,
		Tag:   []byte(""),
		Force: true,
	}
	pcallback := &kinetic.GenericCallback{}
	ph := kinetic.NewResponseHandler(pcallback)
	err = conn.Put(&pentry, ph)
	if err != nil {
		fmt.Println("NonBlocking Put Failure")
	}
	conn.Listen(ph)

	// GET back the object
	gcallback := &kinetic.GetCallback{}
	gh := kinetic.NewResponseHandler(gcallback)
	err = conn.Get(pentry.Key, gh)
	if err != nil {
		fmt.Println("NonBlocking Get Failure")
	}
	conn.Listen(gh)
	gentry := gcallback.Entry

	// Verify the object Key and Value
	if !bytes.Equal(pentry.Key, gentry.Key) {
		fmt.Printf("Key Mismatch: [%s] vs [%s]\n", pentry.Key, gentry.Key)
	}
	if !bytes.Equal(pentry.Value, gentry.Value) {
		fmt.Printf("Value Mismatch: [%s] vs [%s]\n", pentry.Value, gentry.Value)
	}

	// DELETE the object
	dcallback := &kinetic.GenericCallback{}
	dh := kinetic.NewResponseHandler(dcallback)
	dentry := kinetic.Record{
		Key:   pentry.Key,
		Sync:  pentry.Sync,
		Force: true,
	}
	err = conn.Delete(&dentry, dh)
	if err != nil {
		fmt.Println("NonBlocking Delete Failure")
	}
	conn.Listen(dh)
}
