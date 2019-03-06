package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"shared"
)

func createFile(createFileArgs []string) {
	nameNodeAddr, filename, s3Url := parseCreateFileArgs(createFileArgs)

	fileData := downloadFile(s3Url)

	fileSize := fmt.Sprintf("%d", len(fileData))
	createFileResponse := createFileInNameNode(nameNodeAddr, filename, fileSize)

	blocks := makeBlocks(fileData)
	storeAllBlocks(createFileResponse, blocks)
}

func parseCreateFileArgs(args []string) (nameNodeAddr, filename, s3Url string) {
	verboseMessage := fmt.Sprintf("create file with args: %v", args)
	shared.VerbosePrintln(verboseMessage)

	if len(args) != 3 {
		log.Fatal("Input Error: Must use create-file in the following format 'CLI create-file <name-node-address> <filename> <s3-url>")
	}

	nameNodeAddr = args[0]
	filename = args[1]
	s3Url = args[2]

	return
}

func downloadFile(url string) []byte {
	shared.VerbosePrintln("Downloading file from S3 bucket")

	res, err := http.Get(url)
	shared.CheckErrorAndFatal("Unable to download file from S3 bucket URL", err)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Received status code %s when downloading S3 bucket file", res.Status)
	}

	if res.Body == nil {
		log.Fatal("S3 response has no body")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		shared.CheckErrorAndFatal("Unable to read bytes from S3 response body", err)
	}

	shared.VerbosePrintln("Successfully downloaded file")
	return data
}

func createFileInNameNode(nameNodeAddr, filename, size string) (createFileResponse shared.CreateFileNameNodeResponse) {
	shared.VerbosePrintln("Attempting to create file on name node")

	createFileRequest := shared.CreateFileNameNodeRequest{
		FileName: filename,
		Size:     size,
	}
	sendRequestToNameNode(nameNodeAddr, "createFile", createFileRequest, &createFileResponse)

	shared.VerbosePrintln("Successfully created a file on the name node")

	return
}

func makeBlocks(fileData []byte) []string {
	blocks := []string{}
	byteIndex := 0

	for byteIndex < len(fileData) {
		bytesLeftCount := len(fileData) - byteIndex
		endIndex := byteIndex + min(blockSize, bytesLeftCount)

		base64Encoded := base64.StdEncoding.EncodeToString(fileData[byteIndex:endIndex])
		blocks = append(blocks, base64Encoded)

		byteIndex += blockSize
	}

	return blocks
}

func storeAllBlocks(createFileResponse shared.CreateFileNameNodeResponse, blocks []string) {
	if len(createFileResponse.BlockInfos) != len(blocks) {
		log.Fatalf("Name node block list count '%d' does not match calculated blocks count '%d'", len(createFileResponse.BlockInfos), len(blocks))
	}

	shared.VerbosePrintln("Attempting to store all blocks")

	for i, blockInfo := range createFileResponse.BlockInfos {
		storeBlockReq := shared.MakeStoreBlockRequest(blocks[i], blockInfo)
		successful := shared.StoreSingleBlock(storeBlockReq)

		shared.VerbosePrintln(fmt.Sprintf("Attemping to save block (%d/%d)", i, len(blocks)))

		if !successful {
			log.Fatalf("Unable to store block '%s' on any data node", blockInfo.BlockId)
		}
	}

	shared.VerbosePrintln("Successfully stored all blocks to a data node")
}
