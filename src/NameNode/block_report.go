package main

import (
	"encoding/json"
	"net/http"
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

	for i := 0; i < len(myReq.BlockId); i++ {

	}
}
