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
	blockReportTimer = 30 //in seconds
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

	nameNodeUrl := "http://" + nameNodeAddress + shared.PathBlockReport
	if err != nil {
		log.Println(err)
		return
	}
	// grab blocks in directory
	for _, b := range blocks {
		myBlockReport.BlockIds = append(myBlockReport.BlockIds, b.Name())
	}
	// send to NN, since block report we don't needs a response
	buffer, err := shared.ConvertObjectToJsonBuffer(myBlockReport)

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPut, nameNodeUrl, buffer)
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}

	// handle response
	reportResponse  := shared.BlockReportResponse{}
	err = shared.ObjectFromResponse(res, &reportResponse)
	if err != nil {
		log.Println("Error parsing block report", err)
	}
	if reportResponse.Err != "" {
		log.Println(reportResponse.Err)
	}

	return
}

