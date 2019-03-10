package main

/*
call heartbeat from NN

 */

import (
	"log"
	"net/http"
	"os"
)

var (
	selfAddress string
	nameNodeAddress string
	directory string // for testing purposes, change later!
)

func main() {
	/*
	ensure directory exists
	 */
	addr := ""
	if len(os.Args) >= 1 {
		addr = os.Args[1]
	}

	if addr == "" {
		//ipAddress, _ :=  http.Get("http://169.254.169.254/latest/meta-data/public-ipv4)")
		//body, _ := ioutil.ReadAll(ipAddress.Body)
	}

	// accessible var so other store_and_forward can check if self in DnList
	selfAddress = addr

	nameNodeAddress = os.Args[2]

	directory = os.Args[3]
	// os.Args[2] == nameNodeAddress

	/*
	spin up heartbeat and blockreport in separate threads
	send heartbeatreq to NN with IP

	 */

	http.HandleFunc("/getBlock", get_block)
	http.HandleFunc("/storeBlock", store_and_foward)
	http.HandleFunc("/replicateBlocks", replicate_blocks)

	go send_heartbeat()
	go send_block_report()

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panic(err)
	}

}

