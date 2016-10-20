/*
This example demonstrates:

1. Make non-blocking connection to Kinetic drive
2. Start 3000 go routine
3. Each go routine will PUT an object (WRITEBACK mode), Read back the object, then verify content
4. After all go routine done, request FLUSH of all data to persistent media

*/

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	kinetic "github.com/yongzhy/kinetic-go"
)

var option = kinetic.ClientOptions{
	Host: "127.0.0.1", // Test with Simulator
	Port: 8123,
	User: 1,
	Hmac: []byte("asdfasdf")}

func PutGetThread(conn *kinetic.NonBlockConnection, prefix []byte, tid int, rd int, done chan bool) {
	key := []byte(fmt.Sprintf("%s-%05d", prefix, tid))
	v := bytes.Repeat(key, rd)
	if len(v) > 1024*1024 {
		v = v[:1024*1024]
	}

	fmt.Printf("Thread [%02d] Starting %s\n", tid, key)

	pentry := kinetic.Record{
		Key:   key,
		Value: v,
		Sync:  kinetic.SYNC_WRITEBACK,
		Algo:  kinetic.ALGO_SHA1,
		Tag:   []byte(""),
		Force: true,
	}
	pcallback := &kinetic.GenericCallback{}
	ph := kinetic.NewResponseHandler(pcallback)
	err := conn.Put(&pentry, ph)
	if err != nil {
		fmt.Printf("Thread [%02d] PUT FAILURE %s\n", tid, key, err)
		done <- false
		return
	}

	conn.Listen(ph)
	fmt.Printf("Thread [%02d] PUT DONE %s\n", tid, key)

	gcallback := &kinetic.GetCallback{}
	gh := kinetic.NewResponseHandler(gcallback)

	err = conn.Get(key, gh)
	if err != nil {
		fmt.Printf("Thread [%02d] GET FAILURE %s, %s\n", tid, key, err)
		done <- false
		return
	}

	conn.Listen(gh)
	fmt.Printf("Thread [%02d] GET DONE %s\n", tid, key)

	gentry := gcallback.Entry
	if !bytes.Equal(pentry.Key, gentry.Key) {
		fmt.Printf("Thread [%02d] Key Mismatch: [%s] vs [%s]\n", tid, pentry.Key, gentry.Key)
	}
	if !bytes.Equal(pentry.Value, gentry.Value) {
		fmt.Printf("Thread [%02d] Value Mismatch: [%s] vs [%s]\n", tid, pentry.Value, gentry.Value)
	}

	fmt.Printf("Thread [%02d] ALL Done %s\n", tid, key)
	done <- true
}

func main() {
	kinetic.SetLogLevel(kinetic.LogLevelWarn)
	conn, err := kinetic.NewNonBlockConnection(option)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	done := make(chan bool)
	total := 3000

	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= total; i++ {
		go PutGetThread(conn, []byte("Object"), i, rand.Intn(total), done)
	}

	// Wait for all go routine to done
	for i := 1; i <= total; i++ {
		<-done
	}

	// Issue flush to device to flush all data to media
	callback := &kinetic.GenericCallback{}
	h := kinetic.NewResponseHandler(callback)
	err = conn.Flush(h)
	if err != nil {
		fmt.Printf("Flush FAILURE : %s\n", err)
	}
}
