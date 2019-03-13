package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"shared"
)

func createFile(command string, args []string) {
	nameNodeAddr, filename, s3Url := parseCreateFileArgs(command, args)

	file, size := downloadS3FileInFile(s3Url)
	defer os.Remove(tempS3DownloadFileName)

	fileSize := fmt.Sprintf("%d", size)
	createFileResponse := createFileInNameNode(nameNodeAddr, filename, fileSize)

	storeAllBlocks(createFileResponse, file, size)
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
	sendRequestToNameNode(nameNodeAddr, shared.PathFile, http.MethodPut, createFileRequest, &createFileResponse)

	if createFileResponse.Err != "" {
		log.Fatalln(createFileResponse.Err)
	}

	shared.VerbosePrintln("Successfully created a file on the name node")

	return
}

func storeAllBlocks(createFileResponse shared.CreateFileNameNodeResponse, file *os.File, size int64) {
	blocksCount := int(size / blockSize)
	lastBlockSize := size % blockSize
	if lastBlockSize != 0 {
		blocksCount += 1
	}

	if blocksCount != len(createFileResponse.BlockInfos) {
		log.Fatalf("Calculated blocks count '%d' does not match block count in Name Node response '%d'\n", blocksCount, len(createFileResponse.BlockInfos))
	}

	buffer := make([]byte, blockSize)
	for i := 0; i < blocksCount; i++ {
		blockInfo := createFileResponse.BlockInfos[i]
		bytesRead, err := file.Read(buffer)

		shared.VerbosePrintln(fmt.Sprintf("%d. Uploading block '%s'", i, blockInfo.BlockId))

		if err != nil {
			if err != io.EOF {
				shared.CheckErrorAndFatal("Error reading file from s3", err)
			}

			break
		}

		base64EncodedBlock := base64.StdEncoding.EncodeToString(buffer[:bytesRead])
		storeBlockRequest := shared.MakeStoreBlockRequest(base64EncodedBlock, blockInfo)

		shared.StoreSingleBlock(storeBlockRequest)
	}
}
