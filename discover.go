/*
This example demonstrates:

1. Discover all Kinetic drives on subnet defined by CIDR
2. Save live Kinetic drives IP list to output file

*/

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"

	kinetic "github.com/Kinetic/kinetic-go"
)

var cidr = flag.String("cidr", "", "CIDR notion IP address and mask, like \"192.0.2.0/24\" or \"2001:db8::/32\"")
var timeout = flag.Int64("timeout", 50, "Network timeout in milliseconds")
var output = flag.String("output", "drives.txt", "Output file for Kinetic drives IP list")

var wg sync.WaitGroup
var finalIP = make([]string, 0)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func discover(jobs []net.IP, m *sync.Mutex, id int) {
	defer wg.Done()

	for _, ip := range jobs {
		var option = kinetic.ClientOptions{
			Host:    ip.String(),
			Port:    8123,
			User:    1,
			Hmac:    []byte("asdfasdf"),
			Timeout: 50,
		}

		conn, err := kinetic.NewNonBlockConnection(option)
		if err == nil && conn != nil {
			m.Lock()
			finalIP = append(finalIP, ip.String())
			m.Unlock()
		}
	}
}

func main() {
	flag.Parse()

	if cidr == nil || len(*cidr) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	// Set the log leverl to debug
	kinetic.SetLogLevel(kinetic.LogLevelInfo)

	// Parse passed in CIDR IP
	ip, ipnet, err := net.ParseCIDR(*cidr)
	if err != nil {
		flag.Usage()
		panic("Invalid CIDR for IP address and mask")
	}

	// Get total number of IPs to check in the subnet
	var count uint32
	m := []byte(ipnet.Mask)
	count |= uint32(m[0]) << 24
	count |= uint32(m[1]) << 16
	count |= uint32(m[2]) << 8
	count |= uint32(m[3])
	count = count ^ 0xFFFFFFFF
	count++

	// Get all IP address in the subnet
	var allIP = make([]net.IP, 0, count)
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		temp := make([]byte, len(ip))
		copy(temp, ip)
		allIP = append(allIP, temp)
	}

	var ipMu = &sync.Mutex{}
	var worker = (len(allIP) + runtime.NumCPU() - 1) / runtime.NumCPU()
	if worker > 20 {
		worker = 20
	}
	jobCount := len(allIP) / worker
	for w := 0; w < worker; w++ {
		// Spawn worker goroutine to discover given IP range
		go discover(allIP[w*jobCount:(w+1)*jobCount], ipMu, w)
		wg.Add(1)
	}

	// Wait worker to finish
	wg.Wait()

	fmt.Println("################## LIVE KINETIC DRIVES ############################")
	for _, ip := range finalIP {
		fmt.Println(ip)
	}

	// Save result into file
	file, err := os.Create(*output)
	if err != nil {
		fmt.Errorf("Error open file %s to write IP list\n", *output)
	}
	defer file.Close()

	file.WriteString(strings.Join(finalIP, "\n"))
	fmt.Printf("Result wrote into file %s\n", *output)
}
