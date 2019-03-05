package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const (
)

func getFile(write http.ResponseWriter, req *http.Request) { //return dataNode list per block
	decoder := json.NewDecoder(req.Body)
	myReq := createRequest{}
	err := decoder.Decode(&myReq)
	log.Println(err)

	found := false
	fileIndex := -1
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
				found = true
				fileIndex = i
			}
		}
	}

	if !found {
		myRes := responseObj{}
		js, err := convertObjectToJson(myRes)
		log.Print(err)
		write.Header().Set("Content-Type", "application/json")
		write.Write(js)
		return
	}

	myRes := responseObj{}
	myRes.Blocks = make([]createResponse, files.metaData[fileIndex].numBlocks)
	for i := 0; i < files.metaData[fileIndex].numBlocks; i++ {
		blockList := createResponse{}
		blockList.BlockId = myReq.Filename + "_" + strconv.Itoa(i)
		blockList.DnList = make([]string, repFact)
		for j := 0; j < repFact; j++ {
			blockList.DnList[j] = files.metaData[fileIndex].blockLists[i].dnList[j]
		}
		myRes.Blocks[i] = blockList
	}

	js, err := convertObjectToJson(myRes)
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}
