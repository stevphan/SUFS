package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"shared"
	"time"
)

const (
	blockReportTimer = 60 //in seconds
)

func send_block_report() { // sends a block report every blockReportTimer seconds, so every minute
	create_report()
	ticker := time.NewTicker(blockReportTimer * time.Second)
	for range ticker.C {
		log.Println("sending block report")
		create_report()
	}
}

func create_report() {
	myBlockReport := shared.BlockReportRequest{MyIp:selfAddress}

	if !exists(directory) {
		_ = os.MkdirAll(directory, os.ModePerm)
	}

	blocks, err := ioutil.ReadDir(directory)

	nameNodeUrl := "http://" + nameNodeAddress + "/blockReport"
	if err != nil {
		log.Fatal(err)
	}
	// grab blocks in directory
	for _, b := range blocks {
		myBlockReport.BlockIds = append(myBlockReport.BlockIds, b.Name())
	}
	// send as POST to NN, since block report we don't needs a response
	buffer, err := shared.ConvertObjectToJsonBuffer(myBlockReport)
	res, err := http.Post(nameNodeUrl,"application/json", buffer)

	// handle response
	reportResponse  := shared.BlockReportResponse{}
	err = shared.ObjectFromResponse(res, &reportResponse)
	shared.CheckErrorAndFatal("Error sending heartbeat", err)

	return
}

