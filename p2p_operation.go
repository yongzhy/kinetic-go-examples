/*
This example demonstrates:

P2Push operation, at least 2 Kinetic device / simulator needed to run this example.

*/

package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	kinetic "github.com/Kinetic/kinetic-go"
)

var hostName01 = "127.0.0.1" // Simulator 1
var hostName02 = "127.0.0.2" // Simulator 2

var drive01 = kinetic.ClientOptions{
	Host: hostName01,
	Port: 8123,
	User: 1,
	Hmac: []byte("asdfasdf"),
}

var drive02 = kinetic.ClientOptions{
	Host: hostName02,
	Port: 8123,
	User: 1,
	Hmac: []byte("asdfasdf"),
}

func main() {
	kinetic.SetLogLevel(kinetic.LogLevelDebug)

	conn, err := kinetic.NewBlockConnection(drive01)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Step 1 : PUT 2 objects to drive01
	key01 := []byte("P2PObject01")
	entry01 := kinetic.Record{
		Key:   key01,
		Value: []byte(time.Now().Format(time.RFC3339Nano)),
		Sync:  kinetic.SYNC_WRITETHROUGH,
		Algo:  kinetic.ALGO_SHA1,
		Tag:   []byte(""),
		Force: true,
	}
	status, err := conn.Put(&entry01)
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Put 01 failure, ", err, status)
		os.Exit(-1)
	}

	key02 := []byte("P2PObject02")
	entry02 := kinetic.Record{
		Key:   key02,
		Value: []byte(time.Now().Format(time.RFC3339Nano)),
		Sync:  kinetic.SYNC_WRITETHROUGH,
		Algo:  kinetic.ALGO_SHA1,
		Tag:   []byte(""),
		Force: true,
	}
	status, err = conn.Put(&entry02)
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Put 01 failure, ", err, status)
		os.Exit(-1)
	}

	// Step 2 : P2P Operation to PUT object from drive01 to drive 02, under same key
	p2pop := kinetic.P2PPushRequest{
		HostName: hostName02,
		Port:     8123,
		Tls:      false,
		Operations: []kinetic.P2PPushOperation{
			kinetic.P2PPushOperation{
				Key:   key01,
				Force: true,
			},
			kinetic.P2PPushOperation{
				Key:   key02,
				Force: true,
			},
		},
	}
	p2pStatus, status, err := conn.P2PPush(&p2pop)
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("P2PPush failure: ", err, status)
		os.Exit(-1)
	}

	if p2pStatus != nil {
		if p2pStatus.AllOperationsSucceeded == true {
			fmt.Println("P2PPush All Child Operation Succeeded.")
		} else {
			for _, s := range p2pStatus.PushStatus {
				fmt.Println("P2PPush Status : ", s.String())
			}
		}
	}

	// Step 3 : Read objects from drive02 and verify content
	conn2, err := kinetic.NewBlockConnection(drive02)
	if err != nil {
		panic(err)
	}
	defer conn2.Close()

	r01, status, err := conn2.Get(key01)
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Get 01 failure: ", err, status)
		os.Exit(-1)
	}

	r02, status, err := conn2.Get(key02)
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("Get 02 failure: ", err, status)
		os.Exit(-1)
	}

	if !bytes.Equal(r01.Key, key01) || !bytes.Equal(r01.Value, entry01.Value) {
		fmt.Println("ERROR verify object 01 key and content")
		fmt.Println("\tExpect Key: ", key01)
		fmt.Println("\tActual Key: ", r01.Key)
		fmt.Println("\tExpect Value: ", entry01.Value)
		fmt.Println("\tActual Value: ", r01.Value)
	}
	if !bytes.Equal(r02.Key, key02) || !bytes.Equal(r02.Value, entry02.Value) {
		fmt.Println("ERROR verify object 02 key and content")
		fmt.Println("\tExpect Key: ", key02)
		fmt.Println("\tActual Key: ", r02.Key)
		fmt.Println("\tExpect Value: ", entry02.Value)
		fmt.Println("\tActual Value: ", r02.Value)
	}
}
