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

/*
BlockId
 */


 // blocks in directory
 // blocks unique so check filename

type getRequest struct {
	BlockId string `json:"BlockId"`
}

type getResponse struct {
	Block string `json:"Block"`
	Error string `json:"Error"`
}

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
		log.Fatal("Decoding error: ", err)
	}

	fmt.Printf("Received: %s\n", blockReq)

	tempPath := directory + blockReq.BlockId

	returnData := getResponse{
		Block: "",
		Error: "",
	}



	if exists(tempPath) {
		fmt.Println("found " + blockReq.BlockId)
		file, _ := os.Open(tempPath)
		reader := bufio.NewReader(file)
		content, _ := ioutil.ReadAll(reader)
		// encode base64
		returnData.Block = base64.StdEncoding.EncodeToString(content)
		////for testing, print encoded values
		//fmt.Println("ENCODED: " + returnData.Block)
		//
		//// check if decode works by testing decoded value
		//decoded, err := base64.StdEncoding.DecodeString(returnData.Block)
		//if (err != nil) {}
		//
		//// testing, print decoded values (expected: asdf)
		//fmt.Println("decoded: " + string(decoded))
	} else {
		returnData.Error = "404"
	}

	js, err := convertObjectToJson(returnData)
	log.Print(err)
	write.Header().Set("Content-Type", "application/json")
	_, _ = write.Write(js)
	return
	}