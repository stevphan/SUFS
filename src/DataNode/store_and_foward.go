package main

import (
	"encoding/base64"
	"encoding/json"
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



func isError(err error) bool {
	if err != nil {
		log.Println(err.Error())
	}

	return err != nil
}

func removeSelf(storeReq shared.StoreBlockRequest) {
	dnIndex := -1
	found := false
	i := 0
	for i < len(storeReq.DnList) && !found {
		if storeReq.DnList[i] == selfAddress {
			found = true
			dnIndex = i
		}
		i++
	}

	if found {
		storeReq.DnList[dnIndex] = storeReq.DnList[len(storeReq.DnList)-1]
		storeReq.DnList[len(storeReq.DnList)-1] = ""
		storeReq.DnList = storeReq.DnList[:len(storeReq.DnList)-1]
	}
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

	log.Println("==> done creating file", path)
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

	log.Println("==> done writing to file")
}

func store_and_foward(write http.ResponseWriter, req *http.Request)  { // stores Block data into current DataNode, and forwards it to the next DataNode on the list
	defer req.Body.Close()
	storeReq := shared.StoreBlockRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&storeReq)
	if err != nil {
		log.Println("Decoding error: ", err)
	}

	path := directory + "/" + storeReq.BlockId

	// does blocks exists first and foremost? if not, create /blocks/
	if !exists(directory) {
		_ = os.MkdirAll(directory, os.ModePerm)
	}

	// if block is already contained, then find self (if self in DnList) and then forward
	if exists(path) {
		removeSelf(storeReq)
		shared.StoreSingleBlock(storeReq)

		//done
		storeResp := shared.StoreBlockResponse{}
		js, _ := convertObjectToJson(storeResp)
		write.Header().Set("Content-Type", "application/json")
		_, err = write.Write(js)
		return
	}
	// if list is empty, then just stop
	if len(storeReq.DnList) < 1 {
		//done
		storeResp := shared.StoreBlockResponse{}
		js, _ := convertObjectToJson(storeResp)
		write.Header().Set("Content-Type", "application/json")
		_, err = write.Write(js)
		return
	}

	// check if self is in DnList, if not then abort
	var isContained = false
	for _, v := range storeReq.DnList {
		if v == selfAddress {
			isContained = true
		}
	}
	if !isContained {
		returnError(write, "DN_NOT_IN_LIST")
	}

	// decode block data
	log.Println("Block not in Data Node! Saving...")
	decoded, err := base64.StdEncoding.DecodeString(storeReq.Block)

	if err != nil {
		log.Println("Encoding error: ", err)
	}

	// create & write data
	createFile(path)
	writeFile(path, string(decoded))
	// find self and drop from array
	removeSelf(storeReq)
	// forward without self in DnList
	shared.StoreSingleBlock(storeReq)

	storeResp := shared.StoreBlockResponse{}
	js, err := convertObjectToJson(storeResp)
	write.Header().Set("Content-Type", "application/json")
	_, err = write.Write(js)
}

