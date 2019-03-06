package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"shared"
)

func getFile(getFileArgs []string) {
	nameNodeAddr, filename, saveLocation := parseGetFileArgs(getFileArgs)

	getFileResponse := getFileInNameNode(nameNodeAddr, filename)

	blocks := getBlocks(getFileResponse)
	data := reconstructBlocks(blocks)
	saveFileDataToDisk(data, saveLocation)
}

func parseGetFileArgs(args []string) (nameNodeAddr, filename, saveLocation string) {
	verboseMessage := fmt.Sprintf("get file with args: %v", args)
	shared.VerbosePrintln(verboseMessage)

	if len(args) != 3 {
		log.Fatal("Input Error: Must use get-file in the following format 'CLI get-file <name-node-address> <filename> <save-location>")
	}

	nameNodeAddr = args[0]
	filename = args[1]
	saveLocation = args[2]

	return
}

func getFileInNameNode(nameNodeAddr, filename string) (getFileResponse shared.GetFileNameNodeResponse) {
	shared.VerbosePrintln("Attempting to get file from name node")

	getFileRequest := shared.GetFileNameNodeRequest{
		FileName: filename,
	}
	sendRequestToNameNode(nameNodeAddr, "get-file", getFileRequest, &getFileResponse)

	shared.VerbosePrintln("Successfully got a file from the name node")

	return
}

func getBlocks(getFileResponse shared.GetFileNameNodeResponse) (blocks []string) {
	if getFileResponse.Err != "" {
		log.Fatal(getFileResponse.Err)
	}

	blocks = []string{}

	for _, info := range getFileResponse.BlockInfos {
		blocks = append(blocks, getSingleBlock(info))
	}

	return
}

func getSingleBlock(info shared.BlockInfo) string {
	shared.VerbosePrintln("Attempting to get block from data node")

	// for _, dn := range info.DnList {
	// 	getBlockRequest := getBlockRequest{
	// 		BlockId: info.BlockId,
	// 	}

	// 	// TODO: try each DN until one succeeds

	// 	// blockResponse := getBlockResponse{}
	// 	// sendRequestToNode(dn, "get-block", getBlockRequest, &blockResponse)

	// 	// if blockResponse.Err != "" {

	// 	// }
	// }

	shared.VerbosePrintln("Successfully got block from a data node")

	return ""
}

// func getSingleBlockFromDataNode() string {
// 	buffer, err := convertObjectToJsonBuffer(request)
// 	checkErrorAndFatal("Error while communicating with the node:", err)

// 	url := "http://" + nodeAddr + path
// 	res, err := http.Post(url, "application/json", buffer)
// 	checkErrorAndFatal("Error while communicating with the node:", err)

// 	err = objectFromResponse(res, response)
// 	checkErrorAndFatal("Unable to parse response", err)
// }

func reconstructBlocks(blocks []string) (data []byte) {
	data = []byte{}

	for i, block := range blocks {
		decoded, err := base64.StdEncoding.DecodeString(block)
		if err != nil {
			log.Fatalf("Unable to decode block %d/%d: %v", i, len(blocks), err)
		}

		data = append(data, decoded...)
	}

	return
}

func saveFileDataToDisk(data []byte, saveLocation string) {
	file, err := os.Create(saveLocation)
	if err != nil {
		log.Fatal("Unable to create local file:", err)
	}
	defer file.Close()

	written, err := file.Write(data)
	if err != nil {
		log.Fatal("Unable to write data to file:", err)
	}

	if written != len(data) {
		log.Fatalf("Only wrote %d/%d bytes", written, len(data))
	}
}
