package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func getFile(write http.ResponseWriter, req *http.Request) { //return dataNode list per block
	decoder := json.NewDecoder(req.Body)
	myReq := createRequest{}
	err := decoder.Decode(&myReq)
	errorPrint(err)

	//Finds if the file exist
	found := false
	fileIndex := -1 //Index of the file found
	if files.NumFiles < 1 {
		//TODO fail write (DONE?)
		myRes := responseObj{}
		js, err := convertObjectToJson(myRes)
		errorPrint(err)
		write.Header().Set("Content-Type", "application/json")
		_, err = write.Write(js)
		errorPrint(err)
		return
	} else {
		for i := 0; i < files.NumFiles; i++ {
			if files.MetaData[i].FileName == myReq.Filename {
				found = true
				fileIndex = i
			}
		}
	}

	if !found {
		myRes := responseObj{}
		js, err := convertObjectToJson(myRes)
		errorPrint(err)
		write.Header().Set("Content-Type", "application/json")
		_, err = write.Write(js)
		errorPrint(err)
		return
	}

	//Gets the blocks and DnList for the file
	myRes := responseObj{}
	myRes.Blocks = make([]createResponse, files.MetaData[fileIndex].NumBlocks)
	for i := 0; i < files.MetaData[fileIndex].NumBlocks; i++ {
		blockList := createResponse{}
		blockList.BlockId = myReq.Filename + "_" + strconv.Itoa(i)
		blockList.DnList = make([]string, repFact)
		for j := 0; j < repFact; j++ {
			blockList.DnList[j] = files.MetaData[fileIndex].BlockLists[i].DnList[j]
		}
		myRes.Blocks[i] = blockList
	}

	//Returns myRes which is a responseObj
	js, err := convertObjectToJson(myRes)
	errorPrint(err)
	write.Header().Set("Content-Type", "application/json")
	_, err = write.Write(js)
	errorPrint(err)
}
