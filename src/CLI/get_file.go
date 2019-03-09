package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"shared"
	"strings"
)

func getFile(getFileArgs []string, displayDataNodeInfoOnly bool) {
	nameNodeAddr, filename, saveLocation := parseGetFileArgs(getFileArgs, displayDataNodeInfoOnly)

	getFileResponse := getFileInNameNode(nameNodeAddr, filename)

	if getFileResponse.Err != "" {
		log.Fatal(getFileResponse.Err)
	}

	if displayDataNodeInfoOnly {
		displayDataNodeInfo(getFileResponse)
	} else {
		file := createSaveLocation(saveLocation)
		getAndSaveBlocks(getFileResponse, file)
		log.Println("Finished")
	}
}

func parseGetFileArgs(args []string, displayDataNodeInfoOnly bool) (nameNodeAddr, filename, saveLocation string) {
	fmtArgs := stringsMap(args, func(s string) string { return "'" + s + "'" })
	verboseMessage := fmt.Sprintf("get file with args: %v", fmtArgs)
	shared.VerbosePrintln(verboseMessage)

	if displayDataNodeInfoOnly {
		if len(args) != 2 {
			log.Fatal("Input Error: Must use list-data-nodes in the following format 'CLI list-data-nodes <name-node-address> <filename>")
		}
	} else {
		if len(args) != 3 {
			log.Fatal("Input Error: Must use get-file in the following format 'CLI get-file <name-node-address> <filename> <save-location>")
		}

		saveLocation = args[2]
	}

	nameNodeAddr = args[0]
	filename = args[1]

	return
}

func getFileInNameNode(nameNodeAddr, filename string) (getFileResponse shared.GetFileNameNodeResponse) {
	shared.VerbosePrintln("Attempting to get file from name node")

	getFileRequest := shared.GetFileNameNodeRequest{
		FileName: filename,
	}
	sendRequestToNameNode(nameNodeAddr, "getFile", getFileRequest, &getFileResponse)

	shared.VerbosePrintln("Successfully got a file from the name node")

	return
}

func displayDataNodeInfo(getFileResponse shared.GetFileNameNodeResponse) {
	for i, blockInfo := range getFileResponse.BlockInfos {
		dnList := strings.Join(blockInfo.DnList, ", ")
		fmt.Printf("%d. '%s' [%s]\n", i, blockInfo.BlockId, dnList)
	}
}

func createSaveLocation(saveLocation string) *os.File {
	file, err := os.Create(saveLocation)
	if err != nil {
		log.Fatal("Unable to create local file:", err)
	}

	return file
}

func getAndSaveBlocks(getFileResponse shared.GetFileNameNodeResponse, file *os.File) {
	defer file.Close()

	for i, info := range getFileResponse.BlockInfos {
		block, success := getSingleBlock(info)
		if !success {
			os.Remove(file.Name())
			log.Fatalf("Unable to get block %d/%d", i, len(getFileResponse.BlockInfos))
		}

		decoded, err := base64.StdEncoding.DecodeString(block)
		if err != nil {
			os.Remove(file.Name())
			log.Fatalf("Unable to decode block %d/%d: %v", i, len(getFileResponse.BlockInfos), err)
		}

		file.Write(decoded)
	}
}

func getSingleBlock(info shared.BlockInfo) (string, bool) {
	shared.VerbosePrintln("Attempting to get block from data node")

	for _, dn := range info.DnList {
		getBlockRequest := shared.GetBlockRequest{
			BlockId: info.BlockId,
		}

		block, success := getSingleBlockFromDataNode(getBlockRequest, dn)
		if success {
			shared.VerbosePrintln("Successfully got block from a data node")

			return block, true
		}
	}

	return "", false
}

func getSingleBlockFromDataNode(request shared.GetBlockRequest, nodeAddr string) (string, bool) {
	buffer, err := shared.ConvertObjectToJsonBuffer(request)
	if err != nil {
		log.Println("Error while communicating with the data node:", err)
		return "", false
	}

	url := "http://" + nodeAddr + "/store-block"
	res, err := http.Post(url, "application/json", buffer)
	if err != nil {
		log.Println("Error while communicating with the data node:", err)
		return "", false
	}

	response := shared.GetBlockResponse{}
	err = shared.ObjectFromResponse(res, &response)
	if err != nil {
		log.Println("Unable to parse response from the data node:", err)
		return "", false
	}

	if response.Err != "" {
		log.Println("Unable to parse response from the data node:", err)
		return "", false
	}

	return response.Block, true
}
