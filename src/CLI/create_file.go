package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type createFileNameNodeRequest struct {
	FileName string `json:"FileName"`
	Size     string `json:"Size"`
}

type createFileNameNodeResponse struct {
	// TODO: finish me
}

func createFile(createFileArgs []string) {
	nameNodeAddr, filename, s3Url := parseCreateFileArgs(createFileArgs)
	tempFilepath := dirTempCreateFiles + filename

	downloadFile(tempFilepath, s3Url)
	// purposefully ignoring file deletion errors
	defer os.Remove(tempFilepath)

	fileInfo, err := os.Stat(tempFilepath)
	checkErrorAndFatal("Unable to get statistics on temporary file", err)

	createFileRequest := createFileNameNodeRequest{
		FileName: filename,
		Size:     fmt.Sprintf("%d", fileInfo.Size()),
	}

	createFileInNameNode(nameNodeAddr, createFileRequest)

	sendBlocks()

	return
}

func parseCreateFileArgs(args []string) (nameNodeAddr, filename, s3Url string) {
	verboseMessage := fmt.Sprintf("create file with args: %v", args)
	verbosePrintln(verboseMessage)

	if len(args) != 3 {
		log.Fatal("Input Error: Must use get-file in the following format 'CLI get-file <name-node-address> <filename> <s3-url>")
	}

	nameNodeAddr = args[0]
	filename = args[1]
	s3Url = args[2]

	return
}

func downloadFile(filepath string, url string) {
	if useLocalFile {
		verbosePrintln("Assuming file is already downloaded")
		return
	}

	verbosePrintln("Downloading file from S3 bucket")

	// Create the file
	out, err := os.Create(filepath)
	checkErrorAndFatal("Unable to create temporary file", err)
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	checkErrorAndFatal("Unable to download file from S3 bucket URL", err)
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received status code %s when downloading S3 bucket file", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	checkErrorAndFatal("Unable to get S3 bucket file from response", err)

	verbosePrintln("Successfully downloaded file")
}

func createFileInNameNode(nameNodeAddr string, request createFileNameNodeRequest) (createFileResponse createFileNameNodeResponse) {
	verbosePrintln("Requesting name node to create the file")

	nameNodeUrl := "http://" + nameNodeAddr + "/create-file"
	buffer, err := convertObjectToJsonBuffer(request)
	checkErrorAndFatal("Error while communicating to the name node:", err)

	res, err := http.Post(nameNodeUrl, "application/json", buffer)
	checkErrorAndFatal("Error while communicating to the name node:", err)

	err = objectFromResponse(res, &createFileResponse)
	checkErrorAndFatal("Unable to parse response", err)

	return
}

func sendBlocks() {
	// TODO: finish me
}
