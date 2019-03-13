package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"shared"
)

func replicate_blocks(write http.ResponseWriter, req *http.Request)  {
	log.Println("Called to replicate")

	recoverReq := shared.ReplicationRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&recoverReq)
	if err != nil {
		log.Println("Decoding error: ", err)
	}


	recoverResp := shared.ReplicationResponse{}

	tempPath := directory + "/" + recoverReq.BlockId
	if exists(tempPath) {
		log.Println("found " + recoverReq.BlockId)
		file, _ := os.Open(tempPath)
		reader := bufio.NewReader(file)
		content, _ := ioutil.ReadAll(reader)
		// encode base64 the block & send a store request (store_and_forward) on DN list
		storeReq := shared.StoreBlockRequest{DnList:recoverReq.DnList, Block:base64.StdEncoding.EncodeToString(content), BlockId:recoverReq.BlockId}
		shared.StoreSingleBlock(storeReq)
	} else {
		recoverResp.Err = "404_FILE_NOT_FOUND"
	}

	js, err := convertObjectToJson(recoverResp)
	write.Header().Set("Content-Type", "application/json")
	_, err = write.Write(js)
}
