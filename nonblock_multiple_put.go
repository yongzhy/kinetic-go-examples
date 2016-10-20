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

	done := make(chan bool)

	prefix := []byte("TestObject")

	// PUT
	// 1st round: main routin PUT object, start new go routine to wait for operation done
	for id := 1; id <= 100; id++ {
		key := []byte(fmt.Sprintf("%s-%05d", prefix, id))
		v := bytes.Repeat(key, id)
		if len(v) > 1024*1024 {
			v = v[:1024*1024]
		}
		pentry := kinetic.Record{
			Key:   key,
			Value: v,
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

		go func() {
			conn.Listen(ph)
			done <- true
		}()
	}

	// PUT
	// 2nd round, start new go routin for each PUT object and wait for operation done
	for id := 101; id <= 200; id++ {
		go func(id int, done chan bool) {
			key := []byte(fmt.Sprintf("%s-%05d", prefix, id))
			v := bytes.Repeat(key, id)
			if len(v) > 1024*1024 {
				v = v[:1024*1024]
			}
			pentry := kinetic.Record{
				Key:   key,
				Value: v,
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
			done <- true
		}(id, done)
	}

	// Total 200 go routine started, wait for all to finish
	for id := 1; id <= 200; id++ {
		<-done
	}

}
