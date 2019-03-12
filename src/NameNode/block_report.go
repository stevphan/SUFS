package main

import (
	"shared"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func blockReport(write http.ResponseWriter, req *http.Request) { //returns nothing, this is what happens when a block report is received
	myRes := shared.BlockReportResponse{}

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	myReq := shared.BlockReportRequest{}
	err := decoder.Decode(&myReq)
	errorPrint(err)
	log.Print("Block report from ", myReq.MyIp, "\n")

	//finds if the DN in the in dnList, if not add it
	found := false
	//if numDn < 1 {
	if len(dnList) < 1 {
		addToDnList(myReq.MyIp)
		found = true
	} else {
		i := 0
		//for !found && i < numDn {
		for !found && i < len(dnList) {
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
	for i := 0; i < len(myReq.BlockIds); i++ {
		ipFound = false

		var blockFileName string
		var blockIndex int
		var j int //index of values
		found := false
		lock.RLock() //lock before for loop
		for key, value := range files.MetaData { //TODO end this when found, better locking
			lock.RUnlock() //unlock for reads
			j = 0
			for j < len(value) && !found {
				if value[j].Id == myReq.BlockIds[i] {
					found = true
					blockFileName = key
					blockIndex = j
				}
				j++
			}
			lock.RLock() //lock for for loop
		}
		lock.RUnlock()

		if found {
			//tempBlocks, _ := files.MetaData[blockFileName]
			tempBlocks, _ := readMap(blockFileName)
			j = 0
			for j < len(tempBlocks[blockIndex].DnList) && !ipFound {
				if tempBlocks[blockIndex].DnList[j] == myReq.MyIp {
					ipFound = true
				}
				j++
			}
			if !ipFound {
				tempBlocks[blockIndex].DnList = append(tempBlocks[blockIndex].DnList, myReq.MyIp)
				//files.MetaData[blockFileName] = tempBlocks
				writeMap(blockFileName, tempBlocks)
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