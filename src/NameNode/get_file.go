package main

import (
	"encoding/json"
	"log"
	"net/http"
	"Shared"
	"strconv"
)

func getFile(write http.ResponseWriter, req *http.Request) { //return dataNode list per block
	decoder := json.NewDecoder(req.Body)
	myReq := shared.GetFileNameNodeRequest{}
	err := decoder.Decode(&myReq)
	errorPrint(err)

	//Finds if the file exist
	found := false
	fileIndex := -1 //Index of the file found
	//if files.NumFiles < 1 {
	if len(files.MetaData) < 1 {
		myRes := shared.GetFileNameNodeResponse{}
		myRes.Err = "No files found"
		js, err := convertObjectToJson(myRes)
		errorPrint(err)
		write.Header().Set("Content-Type", "application/json")
		_, err = write.Write(js)
		errorPrint(err)
		return
	} else {
		//for i := 0; i < files.NumFiles; i++ {
		for i := 0; i < len(files.MetaData); i++ {
			if files.MetaData[i].FileName == myReq.FileName {
				found = true
				fileIndex = i
			}
		}
	}

	if !found {
		myRes := shared.GetFileNameNodeResponse{}
		myRes.Err = "File " + myReq.FileName + " not found"
		js, err := convertObjectToJson(myRes)
		errorPrint(err)
		write.Header().Set("Content-Type", "application/json")
		_, err = write.Write(js)
		errorPrint(err)
		return
	}

	log.Print("Getting file ", myReq.FileName, "\n")
	//Gets the blocks and DnList for the file
	myRes := shared.GetFileNameNodeResponse{}
	//myRes.BlockInfos = make([]shared.BlockInfo, files.MetaData[fileIndex].NumBlocks)
	myRes.BlockInfos = make([]shared.BlockInfo, len(files.MetaData[fileIndex].BlockLists))
	//for i := 0; i < files.MetaData[fileIndex].NumBlocks; i++ {
	for i := 0; i < len(files.MetaData[fileIndex].BlockLists); i++ {
		blockList := shared.BlockInfo{}
		blockList.BlockId = myReq.FileName + "_" + strconv.Itoa(i)
		blockList.DnList = make([]string, len(files.MetaData[fileIndex].BlockLists[i].DnList))
		for j := 0; j < len(files.MetaData[fileIndex].BlockLists[i].DnList); j++ {
			blockList.DnList[j] = files.MetaData[fileIndex].BlockLists[i].DnList[j]
		}
		myRes.BlockInfos[i] = blockList
	}

	//Returns myRes which is a shared.GetFileNameNodeResponse
	js, err := convertObjectToJson(myRes)
	errorPrint(err)
	write.Header().Set("Content-Type", "application/json")
	_, err = write.Write(js)
	errorPrint(err)
}
