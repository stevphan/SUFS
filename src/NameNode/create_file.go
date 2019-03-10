package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"Shared"
)

const (
	blockSize int64 = 67108864 // assuming bytes, this is equal to 64 MB
)

var (
	currentDn = 0 //index of the DataNode list
)

func createFile(write http.ResponseWriter, req *http.Request) { //needs to return list of dataNodes per block
	var blocksRequired int

	decoder := json.NewDecoder(req.Body)
	myReq := shared.CreateFileNameNodeRequest{}
	err := decoder.Decode(&myReq)
	errorPrint(err)

	//Checks if the file exist
	//if files.NumFiles > 0 {
	if len(files.MetaData) > 0 {
		//for i := 0; i < files.NumFiles; i++ {
		for i := 0; i < len(files.MetaData); i++ {
			if files.MetaData[i].FileName == myReq.FileName {
				myRes := shared.CreateFileNameNodeResponse{}
				myRes.Err = "File with name " + myReq.FileName + " already exist"
				js, err := convertObjectToJson(myRes)
				errorPrint(err)
				write.Header().Set("Content-Type", "application/json")
				_, err = write.Write(js)
				errorPrint(err)
				return
			}
		}
	}

	//Finds the blocks required
	size, err := strconv.ParseInt(myReq.Size, 10, 64)
	errorPrint(err)
	temp := float64(size)/float64(blockSize)
	temp = math.Ceil(temp)
	blocksRequired = int(temp)

	//TODO choose DN to send each block to (check size of, choose lowest) - not super important
	//Checks amount of DN vs the replication factor
	var replicationFactor int
	//if numDn == 0 { //There are no DN
	if len(dnList) == 0 { //There are no DN
		myRes := shared.CreateFileNameNodeResponse{}
		myRes.Err = "No data nodes to store to"
		js, err := convertObjectToJson(myRes)
		errorPrint(err)
		write.Header().Set("Content-Type", "application/json")
		_, err = write.Write(js)
		errorPrint(err)
		return
	//} else if numDn < repFact { //don't have enough DN for replication factor
	} else if len(dnList) < repFact { //don't have enough DN for replication factor
		//replicationFactor = numDn
		replicationFactor = len(dnList)
	} else { //Have enough DN for the replication factor
		replicationFactor = repFact
	}

	//This chooses DN for each block
	//j := 0 //index of the DataNode list
	myRes := shared.CreateFileNameNodeResponse{}
	myRes.BlockInfos = make([]shared.BlockInfo, blocksRequired)
	for i := 0; i < blocksRequired; i++ {
		blockList := shared.BlockInfo{}
		blockList.BlockId = myReq.FileName + "_" + strconv.Itoa(i)
		blockList.DnList = make([]string, replicationFactor)
		for k := 0; k < replicationFactor; k++ {
			blockList.DnList[k] = dnList[currentDn].dnIP
			currentDn++
			//if currentDn == numDn {
			if currentDn == len(dnList) {
				currentDn = 0
			}
		}
		myRes.BlockInfos[i] = blockList
	}

	//Saves the metadata of the file, nothing in DnList yet
	fileToStore := fileMetaData{}
	fileToStore.FileName = myReq.FileName
	//fileToStore.NumBlocks = blocksRequired
	fileToStore.BlockLists = make([]blockList, blocksRequired)
	for i := 0; i < blocksRequired; i++ {
		blockList := blockList{}
		/*for j := 0; j < repFact; j++ {
			blockList.DnListDnList[j] = ""
		}*/
		fileToStore.BlockLists[i] = blockList
	}
	//files.NumFiles++
	files.MetaData = append(files.MetaData, fileToStore)
	writeFilesToDisk()
	
	//Returns myRes which is a shared.CreateFileNameNodeResponse
	js, err := convertObjectToJson(myRes)
	errorPrint(err)
	write.Header().Set("Content-Type", "application/json")
	_, err = write.Write(js)
	errorPrint(err)
	log.Print("File ", myReq.FileName, " created\n")
}