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

func exists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return true
}


func convertObjectToJson(object interface{}) ([]byte, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func get_block(write http.ResponseWriter, req *http.Request) { // returns block requested from the current DN
	blockReq := shared.GetBlockRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&blockReq)
	if err != nil {
		log.Println("Decoding error: ", err)
	}

	log.Print("Received request for block ", blockReq.BlockId, "\n")

	tempPath := directory + "/" + blockReq.BlockId

	returnData := shared.GetBlockResponse{}

	if exists(tempPath) {
		log.Println("found " + blockReq.BlockId)
		file, _ := os.Open(tempPath)
		reader := bufio.NewReader(file)
		content, _ := ioutil.ReadAll(reader)
		// encode base64
		returnData.Block = base64.StdEncoding.EncodeToString(content)
	} else {
		returnData.Err = "Block not found"
	}

	js, err := convertObjectToJson(returnData)
	log.Print(err)
	write.Header().Set("Content-Type", "application/json")
	_, _ = write.Write(js)
	return
}