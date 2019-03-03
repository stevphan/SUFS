package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	blockSize int64 = 67108864 // assuming bytes
)

type createRequest struct{
	Filename string `json:"FileName"`
	Size string		`json:"Size"`
}

type createResponse struct {
	BlockId string	`json:"BlockId"`
	DnList []string `json:"DnList"`
}

type responseObj struct { //What is returned to cli
	Blocks []createResponse
}

func createFile(write http.ResponseWriter, req *http.Request) { //needs to return list of dataNodes per block
	var blocksRequired int

	decoder := json.NewDecoder(req.Body)
	myReq := createRequest{}
	err := decoder.Decode(&myReq)
	log.Println(err)

	//Checks if the file exist
	if files.numFiles < 1 {
		//TODO fail write (DONE?)
		myRes := responseObj{}
		js, err := convertObjectToJson(myRes)
		log.Print(err)
		write.Header().Set("Content-Type", "application/json")
		write.Write(js)
		return
	} else {
		for i := 0; i < files.numFiles; i++ {
			if files.metaData[i].fileName == myReq.Filename {
				//TODO fail write (DONE?)
				myRes := responseObj{}
				js, err := convertObjectToJson(myRes)
				log.Print(err)
				write.Header().Set("Content-Type", "application/json")
				write.Write(js)
				return
			}
		}
	}

	size, err := strconv.ParseInt(myReq.Size, 10, 64)
	if err != nil {
		log.Print("Error with converting string to int64")
		log.Print(err)
	} else {
		blocksRequired = int(size/blockSize)
	}

	fmt.Println(blocksRequired)

	//TODO choose DN to send each block to (check size of, choose lowest)

	var replicationFactor int
	j := 0 //index of the DataNode list
	if numDn == 0 {
		//TODO fail write (DONE?)
		myRes := responseObj{}
		js, err := convertObjectToJson(myRes)
		log.Print(err)
		write.Header().Set("Content-Type", "application/json")
		write.Write(js)
		return
	} else if numDn < repFact {
		replicationFactor = numDn
	} else {
		replicationFactor = repFact
	}

	myRes := responseObj{}
	myRes.Blocks = make([]createResponse, blocksRequired)

	for i := 0; i < blocksRequired; i++ {
		blockList := createResponse{}
		blockList.BlockId = myReq.Filename + "_" + strconv.Itoa(i)
		blockList.DnList = make([]string, replicationFactor)
		for k := 0; k < replicationFactor; k++ {
			blockList.DnList[k] = dnList[j]
			j++
			if j == numDn {
				j = 0
			}
		}
		myRes.Blocks[i] = blockList
	}

	fmt.Println(myRes)
	//TODO return the array of createResponse struct (DONE?)

	//js, err := convertObjectToJsonBuffer(myRes)

	/*js, err := json.Marshal(myRes)
	if err != nil {
		http.Error(write, "Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}*/

	js, err := convertObjectToJson(myRes)
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}