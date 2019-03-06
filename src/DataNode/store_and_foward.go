package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

/*
Block
BlockId
DataNodeList
 */

 // passed in JSON payloads

type storeBlockRequest struct{
	Block string `json:"Block"`
	DnList[] string `json:"DataNodeList"`
	BlockId string `json:"BlockId"`

}


type storeBlockResponse struct {
	Error string `json:"Error"`
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return err != nil
}

func createFile(path string) {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) { return }
		defer file.Close()
	}

	fmt.Println("==> done creating file", path)
}

func writeFile(path string, data string) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) { return }
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(data)
	if isError(err) { return }

	// save changes
	err = file.Sync()
	if isError(err) { return }

	fmt.Println("==> done writing to file")
}

func store_and_foward(write http.ResponseWriter, req *http.Request)  { // stores Block data into current DataNode, and forwards it to the next DataNode on the list
	storeReq := storeBlockRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&storeReq)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	fmt.Printf("Received: %s\n", storeReq)
	fmt.Println("")

	path := s3address + storeReq.BlockId

	// check if file exists already
	if exists(path) { // for debug purposes, otherwise dont print anything
		fmt.Println("FOUND! ... dont do anything")
		return
	} else {
		fmt.Println("Block not in Data Node! Saving...")
		decoded, err := base64.StdEncoding.DecodeString(storeReq.Block)

		if isError(err) { // if there is an error send a response
			errorReq := storeBlockResponse{}
			errorReq.Error = "some_error"
			js, err := convertObjectToJson(errorReq)
			log.Print(err)
			write.Header().Set("Content-Type", "application/json")
			_, _ = write.Write(js)
			return
		}

		fmt.Println("Decoded block data: " + string(decoded))
		createFile(path)
		writeFile(path, string(decoded))
	}



	//	TODO: figure out 'forwarding'
	//	TODO: if block already exists, what about overwriting old data (i.e. adding a line to a file)? should I write anyway?
	// if I have the address/location of each DN, maybe insert directly into the DN??? would have to set up and test if possible [

}
