package main

import (
	"net/http"
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

	for i, info := range getFileResponse.BlockInfos {
		block, success := getSingleBlock(info)
		if !success {
			log.Fatalf("Unable to get block %d/%d", i, len(getFileResponse.BlockInfos))
		}

		blocks = append(blocks, block)
	}

	return
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
