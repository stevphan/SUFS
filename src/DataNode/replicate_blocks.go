package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"shared"
)

func replicate_blocks(write http.ResponseWriter, req *http.Request)  {
	recoverReq := shared.ReplicationRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&recoverReq)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	nameNodeUrl := "http://" + nameNodeAddress + "/blockReport"

	recoverResp := shared.ReplicationResponse{Err:""}

	tempPath := directory + recoverReq.BlockId
	if exists(tempPath) {
		fmt.Println("found " + recoverReq.BlockId)
		file, _ := os.Open(tempPath)
		reader := bufio.NewReader(file)
		content, _ := ioutil.ReadAll(reader)
		// encode base64 the block & send a store request (store_and_forward) on DN list
		storeReq := shared.StoreBlockRequest{DnList:recoverReq.DnList, Block:base64.StdEncoding.EncodeToString(content), BlockId:recoverReq.BlockId}
		shared.StoreSingleBlock(storeReq)
	} else {
		recoverResp.Err = "404_FILE_NOT_FOUND"
	}

	// return POST to
	buffer, err := shared.ConvertObjectToJsonBuffer(recoverResp)
	res, err := http.Post(nameNodeUrl,"application/json", buffer)
	err = shared.ObjectFromResponse(res, &recoverResp)
	shared.CheckErrorAndFatal("Error sending replication", err)

	return
}
