package main

/*
call heartbeat from NN

 */

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"shared"
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
	selfAddress = string(body) + ":8080"
	log.Print("MyIp: ", selfAddress, "\n")

	nameNodeAddress = os.Args[1]
	directory = os.Args[2]
	// os.Args[2] == nameNodeAddress

	/*
	spin up heartbeat and blockreport in separate threads
	send heartbeatreq to NN with IP

	 */

	blockPath := make(map[string]func(http.ResponseWriter, *http.Request))
	blockPath[http.MethodGet] = get_block
	blockPath[http.MethodPut] = store_and_foward
	shared.ServeCall(shared.PathBlock, blockPath)

	//http.HandleFunc("/getBlock", get_block)
	//http.HandleFunc("/storeBlock", store_and_foward)

	repPath := make(map[string]func(http.ResponseWriter, *http.Request))
	blockPath[http.MethodPost] = replicate_blocks
	shared.ServeCall(shared.PathReplication, repPath)

	//http.HandleFunc("/replicateBlocks", replicate_blocks)

	go send_heartbeat()
	go send_block_report()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}

}

