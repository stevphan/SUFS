package main

/*
call heartbeat from NN

 */

import (
	"io/ioutil"
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
	ipAddressRes, _ :=  http.Get("http://169.254.169.254/latest/meta-data/public-ipv4")
	body, _ := ioutil.ReadAll(ipAddressRes.Body)
	selfAddress = string(body)
	log.Print("MyIp: ", selfAddress, "\n")

	nameNodeAddress = os.Args[1]
	directory = os.Args[2]
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

	err := http.ListenAndServe(selfAddress, nil)
	if err != nil {
		log.Panic(err)
	}

}

