package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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


func convertObjectToJsonBuffer(object interface{}) (*bytes.Buffer, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(data)

	return buffer, nil
}

func convertObjectToJson(object interface{}) ([]byte, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// DN will have it's own S3 URL to save to, so for now just save to a folder on disk

func get_block(write http.ResponseWriter, req *http.Request) { // returns block requested from the current DN
	blockId := getRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&blockId)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	fmt.Printf("Received: %s\n", blockId)

	tempPath := "/Users/stxv/blocks/" + blockId.BlockId

	returnData := &getResponse{
		Block: "",
		Error: "",
	}

	if exists(tempPath) {
		fmt.Println("found " + blockId.BlockId)
		file, _ := os.Open(tempPath)
		reader := bufio.NewReader(file)
		content, _ := ioutil.ReadAll(reader)

		// encode base64
		returnData.Block = base64.StdEncoding.EncodeToString(content)

		// for testing, print encoded values
		//fmt.Println("ENCODED: " + encoded)
		//
		//// check if decode works by testing decoded value
		//decoded, err := base64.StdEncoding.DecodeString(encoded)
		//if (err != nil) {
		//
		//}
		//// testing, print decoded values (expected: asdf)
		//fmt.Println("decoded: " + string(decoded))



	} else {
		fmt.Println("did not find " + blockId.BlockId)
		returnData.Error = "404"
	}

	returnContent, _ := convertObjectToJsonBuffer(returnData)

	return returnContent


}