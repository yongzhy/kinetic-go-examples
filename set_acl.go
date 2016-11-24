/*
This example demonstrates:

1. Connect use user id 1
2. Set user id 100 only have GetLog permission
3. Verify user id permission by: GetLog should be OK, GET will fail
*/

package main

import (
	"fmt"

	kinetic "github.com/Kinetic/kinetic-go"
)

func setACL() {
	// Client options
	var option = kinetic.ClientOptions{
		Host:   "127.0.0.1", // Test with Simulator
		Port:   8443,        // Must be SSL connection here
		User:   1,
		Hmac:   []byte("asdfasdf"),
		UseSSL: true, // Set ACL must use SSL connection
	}

	conn, err := kinetic.NewBlockConnection(option)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	perms := []kinetic.ACLPermission{
		kinetic.ACLPermissionGetLog,
	}
	scope := []kinetic.ACLScope{
		kinetic.ACLScope{
			Permissions: perms,
		},
	}
	acls := []kinetic.ACL{
		kinetic.ACL{
			Identify: 100,
			Key:      []byte("asdfasdf"),
			Algo:     kinetic.ACLAlgorithmHMACSHA1,
			Scopes:   scope,
		},
	}

	status, err := conn.SetACL(acls)
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("SetACL failure: ", err, status)
	}

}

func verifyACL() {
	// Client options
	var option = kinetic.ClientOptions{
		Host: "127.0.0.1",
		Port: 8123,
		User: 100,
		Hmac: []byte("asdfasdf")}

	conn, err := kinetic.NewBlockConnection(option)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	logs := []kinetic.LogType{
		kinetic.LogTypeUtilizations,
		kinetic.LogTypeTemperatures,
		kinetic.LogTypeCapacities,
		kinetic.LogTypeConfiguration,
		kinetic.LogTypeStatistics,
		kinetic.LogTypeMessages,
		kinetic.LogTypeLimits,
	}

	_, status, err := conn.GetLog(logs)
	if err != nil || status.Code != kinetic.OK {
		fmt.Println("GetLog Failure: ", err, status)
	}

	_, status, err = conn.Get([]byte("object000"))
	if err != nil {
		fmt.Println("Get Failure: ", err)
	}

	if status.Code != kinetic.RemoteNotAuthorized {
		fmt.Println("SET ACL not effective, ", status)
	}
}

func main() {
	// Set the log leverl to debug
	kinetic.SetLogLevel(kinetic.LogLevelDebug)

	setACL()
	verifyACL()
}
