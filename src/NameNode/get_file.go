package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
)

func getFile(write http.ResponseWriter, req *http.Request) { //return dataNode list per block
	decoder := json.NewDecoder(req.Body)
	myReq := createRequest{}
	err := decoder.Decode(&myReq)
	log.Println(err)

	found := false
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
}
