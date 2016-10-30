/*
This example demonstrates:

Batch PUT objects to kinetic drive.

Batch operation support Batch PUT and DELETE

*/

package main

import (
	"bytes"
	"fmt"

	kinetic "github.com/Kinetic/kinetic-go"
)

var option = kinetic.ClientOptions{
	Host: "127.0.0.1", // Test with simulator
	Port: 8123,
	User: 1,
	Hmac: []byte("asdfasdf")}

func main() {
	kinetic.SetLogLevel(kinetic.LogLevelDebug)

	conn, err := kinetic.NewNonBlockConnection(option)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	prefix := []byte("TestBatchObjects")

	callback := &kinetic.GenericCallback{}
	h := kinetic.NewResponseHandler(callback)
	err = conn.BatchStart(h)
	if err != nil {
		panic("Batch Can't Start " + err.Error())
	}

	err = conn.Listen(h)

	if err != nil || callback.Status().Code != kinetic.OK {
		fmt.Println("Error start batch operation :", err, callback.Status())
		return
	}

	// PUT
	for id := 1; id <= 2; id++ {
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
		err = conn.BatchPut(&pentry)
		if err != nil {
			fmt.Println("NonBlocking Put Failure")
		}
		// Batch PUT / DELETE doesn't have Response message from kinetic device
		// Response message only for in END_BATCH, or ABORT_BATCH
	}

	bcallback := &kinetic.BatchEndCallback{}
	h = kinetic.NewResponseHandler(bcallback)
	err = conn.BatchEnd(h)
	//err = conn.BatchAbort(h)
	if err != nil {
		panic("Batch Can't Start " + err.Error())
	}

	err = conn.Listen(h)
	if err != nil || bcallback.Status().Code != kinetic.OK {
		fmt.Println("Error start batch operation :", err, callback.Status())
		return
	}

	fmt.Println("Batch DONE sequence: ", bcallback.BatchStatus.DoneSequence)
	fmt.Println("Batch FAILED sequence: ", bcallback.BatchStatus.FailedSequence)
}
