package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
BlockId
 */


 // blocks in directory
 // blocks unique so check filename

type getRequest struct {
	BlockId string `json:"BlockId"`
}


func get_block(write http.ResponseWriter, req *http.Request) { // returns block requested from the current DN
	decoder := json.NewDecoder(req.Body)
	myReq := getRequest{}
	err := decoder.Decode(&myReq)

	log.Println(err)
	fmt.Println("")

}