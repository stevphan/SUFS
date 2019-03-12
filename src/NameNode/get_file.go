package main

import (
	"shared"
	"encoding/json"
	"log"
	"net/http"
)

func getFile(write http.ResponseWriter, req *http.Request) { //return dataNode list per block
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	myReq := shared.GetFileNameNodeRequest{}
	err := decoder.Decode(&myReq)
	errorPrint(err)

	//Finds if the file exist
	//_, found := files.MetaData[myReq.FileName]
	_, found := readMap(myReq.FileName)
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

	//Gets the blocks and DnList for the file
	log.Print("Getting file ", myReq.FileName, "\n")
	myRes := shared.GetFileNameNodeResponse{}
	//myBlocks := files.MetaData[myReq.FileName]
	myBlocks, _ := readMap(myReq.FileName)
	myRes.BlockInfos = make([]shared.BlockInfo, len(myBlocks))
	for i := 0; i < len(myBlocks); i++ {
		blockList := shared.BlockInfo{}
		blockList.BlockId = myBlocks[i].Id
		blockList.DnList = myBlocks[i].DnList
		myRes.BlockInfos[i] = blockList
	}

	//Returns myRes which is a shared.GetFileNameNodeResponse
	js, err := convertObjectToJson(myRes)
	errorPrint(err)
	write.Header().Set("Content-Type", "application/json")
	_, err = write.Write(js)
	errorPrint(err)
}
