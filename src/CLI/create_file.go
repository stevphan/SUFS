package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"shared"
)

func createFile(command string, args []string) {
	nameNodeAddr, filename, s3Url := parseCreateFileArgs(command, args)

	fileData := downloadS3File(s3Url)

	fileSize := fmt.Sprintf("%d", len(fileData))
	createFileResponse := createFileInNameNode(nameNodeAddr, filename, fileSize)

	blocks := makeBlocks(fileData)
	storeAllBlocks(createFileResponse, blocks)
}

func parseCreateFileArgs(command string, args []string) (nameNodeAddr, filename, s3Url string) {
	fmtArgs := stringsMap(args, func(s string) string { return "'" + s + "'" })
	verboseMessage := fmt.Sprintf("create file with args: %v", fmtArgs)
	shared.VerbosePrintln(verboseMessage)

	if len(args) != 3 {
		log.Fatalf("Input Error: Must use create-file in the following format '%s create-file <name-node-address> <filename> <s3-url>'\n", command)
	}

	nameNodeAddr = args[0]
	filename = args[1]
	s3Url = args[2]

	return
}

func createFileInNameNode(nameNodeAddr, filename, size string) (createFileResponse shared.CreateFileNameNodeResponse) {
	shared.VerbosePrintln("Attempting to create file on name node")

	createFileRequest := shared.CreateFileNameNodeRequest{
		FileName: filename,
		Size:     size,
	}
	sendRequestToNameNode(nameNodeAddr, "createFile", createFileRequest, &createFileResponse)

	if createFileResponse.Err != "" {
		log.Fatalln(createFileResponse.Err)
	}

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
		log.Fatalf("Name node block list count '%d' does not match calculated blocks count '%d'\n", len(createFileResponse.BlockInfos), len(blocks))
	}

	shared.VerbosePrintln("Attempting to store all blocks")

	for i, blockInfo := range createFileResponse.BlockInfos {
		storeBlockReq := shared.MakeStoreBlockRequest(blocks[i], blockInfo)
		successful := shared.StoreSingleBlock(storeBlockReq)

		shared.VerbosePrintln(fmt.Sprintf("Attemping to save block (%d/%d)", i, len(blocks)))

		if !successful {
			log.Fatalf("Unable to store block '%s' on any data node\n", blockInfo.BlockId)
		}
	}

	shared.VerbosePrintln("Successfully stored all blocks to a data node")
}
