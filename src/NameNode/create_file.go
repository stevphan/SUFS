package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"math"
	"net/http"
)

const (
	blockSize float64 = 67108864 // assuming bytes
)

type createRequest struct{
	Filename string `json:"FileName"`
	Size string		`json:"Size"`
}

type createResponse struct {
	blockId string	`json:"BlockId"`
	dnList []string `json:"DnList"`
}

func createFile(write http.ResponseWriter, req *http.Request) { //needs to return list of dataNodes per block
	var blocksRequired int64

	decoder := json.NewDecoder(req.Body)
	myReq := createRequest{}
	err := decoder.Decode(&myReq)

	log.Println(err)

	//TODO make sure file doesn't exist

	size, err := strconv.ParseInt(myReq.Size, 10, 64)
	if err != nil {
		log.Print("Error with converting string to int64")
		log.Print(err)
	} else {
		blocksRequired = math.Ceil(size/blockSize)
	}

	fmt.Println(blocksRequired)
	fmt.Println()

	//TODO choose DN to send each block to (check size of, choose lowest)
	//
	//TODO return the array of createResponse struct


}
