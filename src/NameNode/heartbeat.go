package main

import (
	"Shared"
	"encoding/json"
	"net/http"
)

func heartBeat(write http.ResponseWriter, req *http.Request) {
	myRes := shared.HeartbeatResponse{}

	decoder := json.NewDecoder(req.Body)
	myReq := shared.HeartbeatRequest{}
	err := decoder.Decode(&myReq)
	errorPrint(err)

	//TODO delay some timer

	//Returns myRes which is a shared.HeartbeatResponse
	js, err := convertObjectToJson(myRes)
	errorPrint(err)
	write.Header().Set("Content-Type", "application/json")
	_, err = write.Write(js)
	errorPrint(err)
}
