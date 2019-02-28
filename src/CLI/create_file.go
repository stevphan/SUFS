package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func createFile(args []string) {
	/*
		0: filename
		1: s3 bucket url
	*/

	// TODO: check for 2 arguments
	// TODO: check for valid filename
	// TODO: check for valid s3 bucket url

	log.Println("create file args: ", args)

	if len(args) != 2 {
		log.Fatal("Input Error: Must use get-file in the following format 'CLI get-file <filename> <s3-url>")

	}
	filepath := args[0]
	s3Url := args[1]

	// filePathBase := "/Users/Rivukis/Desktop/tmp/"
	// url := "https://s3.amazonaws.com/amazon-reviews-pds/tsv/amazon_reviews_us_Electronics_v1_00.tsv.gz"
	// filepath := filePathBase + "mys3file"

	downloadFile(filepath, s3Url)

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("file size", fileInfo.Size())

	return
}

func downloadFile(filepath string, url string) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Non 200 response when downloading file: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
