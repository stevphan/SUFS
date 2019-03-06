package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"shared"
)

/*
Block
BlockId
DataNodeList
 */

 // passed in JSON payloads


func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
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

func returnError(write http.ResponseWriter, errMessage string) {
	 // if there is an error send a response
		errorReq := shared.StoreBlockResponse{}
		errorReq.Err = errMessage
		js, _ := convertObjectToJson(errorReq)
		write.Header().Set("Content-Type", "application/json")
		_, _ = write.Write(js)
		return

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
	storeReq := shared.StoreBlockRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&storeReq)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	fmt.Printf("Received: %s\n", storeReq)
	fmt.Println("")

	path := folderPath + storeReq.BlockId

	// FOR TESTING ... will be IP addr later and wont be able to
	s3address := "test_node_2.aws.com"

	// if list is empty, then just stop
	if len(storeReq.DnList) < 1 {
		return
	}

	// check if self is in DnList, if not then abort
	var isContained = false
	for _, v := range storeReq.DnList {
		if v == s3address {
			isContained = true
		}
	}
	if !isContained {
		returnError(write, "DN_NOT_IN_LIST")
	}

	// decode block data
	fmt.Println("Block not in Data Node! Saving...")
	decoded, err := base64.StdEncoding.DecodeString(storeReq.Block)

	if err != nil {
		log.Fatal("Encoding error: ", err)
	}

	fmt.Println("Decoded block data: " + string(decoded))
	// create & write data
	createFile(path)
	writeFile(path, string(decoded))
	// find self and drop from array
	for i, v := range storeReq.DnList {
		if v == s3address {
			storeReq.DnList = append(storeReq.DnList[:i], storeReq.DnList[i+1:]...)
			break
		}
	}
	// forward without self in DnList
	shared.StoreSingleBlock(storeReq)

	fmt.Println()

	/*
	on success, drop self from array
	if array not empty
	forward to first in array
	 */

	//	TODO: figure out 'forwarding'
	//	TODO: if block already exists, what about overwriting old data (i.e. adding a line to a file)? should I write anyway?
	// if I have the address/location of each DN, maybe insert directly into the DN??? would have to set up and test if possible [

}

