package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type nameNodeRequest struct {
	FileName string `json:"FileName"`
	Size     string `json:"Size"`
}

type nameNodeResponse struct {
	// TODO: finish me
}

func createFile(args []string) {
	msg := fmt.Sprintf("create file with args: %v", args)
	verbosePrintln(msg)

	if len(args) != 3 {
		log.Fatal("Input Error: Must use get-file in the following format 'CLI get-file <name-node-address> <filename> <s3-url>")
	}

	nameNodeAddr := args[0]
	filename := args[1]
	s3Url := args[2]
	tempFilepath := "/Users/Rivukis/Desktop/tmp/" + filename

	if useLocalFile {
		verbosePrintln("Assuming file is already downloaded")
	} else {
		downloadFile(tempFilepath, s3Url)
		defer deleteTempFile(tempFilepath)
	}

	fileInfo, err := os.Stat(tempFilepath)
	checkErrorAndFatal("Unable to get statistics on temporary file", err)

	nnRequest := nameNodeRequest{
		FileName: filename,
		Size:     fmt.Sprintf("%d", fileInfo.Size()),
	}

	createFileInNameNode(nameNodeAddr, nnRequest)

	return
}

func downloadFile(filepath string, url string) {
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

func createFileInNameNode(nameNodeAddr string, request nameNodeRequest) (nnRes nameNodeResponse) {
	verbosePrintln("Requesting name node to create the file")

	nameNodeUrl := "http://" + nameNodeAddr + "/create-file"
	buffer, err := convertObjectToJsonBuffer(request)
	checkErrorAndFatal("Error while communicating to the name node:", err)

	res, err := http.Post(nameNodeUrl, "application/json", buffer)
	checkErrorAndFatal("Error while communicating to the name node:", err)

	err = objectFromResponse(res, &nnRes)
	checkErrorAndFatal("Unable to parse response", err)

	return
}

func deleteTempFile(filepath string) {
	// TODO: finish me
}

func sendBlocks() {
	// TODO: finish me
}
