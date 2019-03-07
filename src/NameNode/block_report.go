package main

import (
	"Shared"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func blockReport(write http.ResponseWriter, req *http.Request) { //returns nothing, this is what happens when a block report is received
	myRes := shared.BlockReportResponse{}

	decoder := json.NewDecoder(req.Body)
	myReq := shared.BlockReportRequest{}
	err := decoder.Decode(&myReq)
	errorPrint(err)

	found := false
	if numDn < 1 {
		dnList = append(dnList, myReq.MyIp)
		numDn++
	} else {
		//for i := 0; i < numDn; i++ {
		i := 0
		for !found && i < numDn {
			if dnList[i] == myReq.MyIp {
				found = true
			}
			i++
		}
	}
	if !found {
		dnList = append(dnList, myReq.MyIp)
		numDn++
	}

	//TODO make sure files have the same number of blocks
	var ipFound bool
	for i := 0; i < len(myReq.BlockIds); i++ {
		ipFound = false
		temp := strings.Split(myReq.BlockIds[i], "_")
		fileName, blockNum := temp[0], temp[1]
		blockId, err := strconv.Atoi(blockNum)
		errorPrint(err)
		found, fileIndex := findFile(fileName)
		if found {
			for j := 0; j < len(files.MetaData[fileIndex].BlockLists[blockId].DnList); j++ {
				//fmt.Println(files)
				if myReq.MyIp == files.MetaData[fileIndex].BlockLists[blockId].DnList[j]{
					ipFound = true
				}
			}
			if !ipFound {
				/*j := 0
				for !ipFound && j < len(files.MetaData[fileIndex].BlockLists[blockId].DnList){
					if files.MetaData[fileIndex].BlockLists[blockId].DnList[j] == "" {
						files.MetaData[fileIndex].BlockLists[blockId].DnList[j] = myReq.MyIp
						ipFound = true
					}
					j++
				}*/
				files.MetaData[fileIndex].BlockLists[blockId].DnList = append(files.MetaData[fileIndex].BlockLists[blockId].DnList, myReq.MyIp)
			}
		}
	}
	//fmt.Println(files)
	writeFilesToDisk()

	//Returns myRes which is a shared.BlockReportResponse
	js, err := convertObjectToJson(myRes)
	errorPrint(err)
	write.Header().Set("Content-Type", "application/json")
	_, err = write.Write(js)
	errorPrint(err)
}