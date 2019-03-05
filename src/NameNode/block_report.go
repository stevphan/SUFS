package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
)

func blockReport(write http.ResponseWriter, req *http.Request) { //returns nothing, this is what happens when a block report is received
	decoder := json.NewDecoder(req.Body)
	myReq := blockReportRequest{}
	err := decoder.Decode(&myReq)
	errorPrint(err)

	found := false
	if numDn < 1 {
		dnList = append(dnList, myReq.MyIp)
		numDn++
	} else {
		for i := 0; i < numDn; i++ {
			if dnList[i] == myReq.MyIp {
				found = true
			}
		}
	}
	if !found {
		dnList = append(dnList, myReq.MyIp)
		numDn++
	}

	var ipFound bool
	for i := 0; i < len(myReq.BlockId); i++ {
		ipFound = false
		temp := strings.Split(myReq.BlockId[i], "_")
		fileName, blockNum := temp[0], temp[1]
		blockId, err := strconv.Atoi(blockNum)
		errorPrint(err)
		found, fileIndex := findFile(fileName)
		if found {
			for j := 0; j < repFact; j++ {
				if myReq.MyIp == files.MetaData[fileIndex].BlockLists[blockId].DnList[j]{
					ipFound = true
				}
			}
			if !ipFound {
				j := 0
				for !ipFound || j < repFact{
					if files.MetaData[fileIndex].BlockLists[blockId].DnList[j] == "" {
						files.MetaData[fileIndex].BlockLists[blockId].DnList[j] = myReq.MyIp
					}
				}
			}
		}
	}
	fmt.Println()
}

//Returns true if file name found plus the index it is in the files.metadata[]
//Returns false if file name not found plus index of -1
func findFile(fileName string) (found bool, fileIndex int) {
	if files.NumFiles > 0 {
		for i := 0; i < files.NumFiles; i++ {
			if files.MetaData[i].FileName == fileName {
				return true, i
			}
		}
	}
	return false, -1
}