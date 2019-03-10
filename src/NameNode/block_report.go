package main

import (
	"Shared"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func blockReport(write http.ResponseWriter, req *http.Request) { //returns nothing, this is what happens when a block report is received
	myRes := shared.BlockReportResponse{}

	decoder := json.NewDecoder(req.Body)
	myReq := shared.BlockReportRequest{}
	err := decoder.Decode(&myReq)
	errorPrint(err)
	log.Print("Block report from ", myReq.MyIp, "\n")

	//finds if the DN in the in dnList, if not add it
	found := false
	if numDn < 1 {
		addToDnList(myReq.MyIp)
		found = true
	} else {
		i := 0
		for !found && i < numDn {
			if dnList[i].dnIP == myReq.MyIp {
				found = true
				dnList[i].dnTime = time.Now() //got a response, so delay "death"
			}
			i++
		}
	}
	if !found {
		addToDnList(myReq.MyIp)
	}

	var ipFound bool
	var j int //index for a block's dnList
	for i := 0; i < len(myReq.BlockIds); i++ {
		ipFound = false
		temp := strings.Split(myReq.BlockIds[i], "_")
		fileName, blockNum := temp[0], temp[1]
		blockId, err := strconv.Atoi(blockNum)
		errorPrint(err)
		found, fileIndex := findFile(fileName)
		if found {
			if blockId < files.MetaData[fileIndex].NumBlocks { //makes sure blockId exist
				//for j := 0; j < len(files.MetaData[fileIndex].BlockLists[blockId].DnList); j++ {
				j = 0
				for j < len(files.MetaData[fileIndex].BlockLists[blockId].DnList) && !ipFound {
					if myReq.MyIp == files.MetaData[fileIndex].BlockLists[blockId].DnList[j]{
						ipFound = true
					}
					j++
				}
				if !ipFound {
					files.MetaData[fileIndex].BlockLists[blockId].DnList = append(files.MetaData[fileIndex].BlockLists[blockId].DnList, myReq.MyIp)
				}
			}
		}
	}
	writeFilesToDisk()

	//Returns myRes which is a shared.BlockReportResponse
	js, err := convertObjectToJson(myRes)
	errorPrint(err)
	write.Header().Set("Content-Type", "application/json")
	_, err = write.Write(js)
	errorPrint(err)
}