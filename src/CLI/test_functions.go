package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"shared"
)

func doRequestTesting() {
	// originalObject := createFileNameNodeRequest{
	// 	FileName: "myfilename",
	// 	Size: "12345",
	// }

	dnList := []string{"1.2.3.4", "5.6.7.8", "10.0.0.7"}

	originalObject := shared.StoreBlockRequest{
		Block:   "actualblock",
		DnList:  dnList,
		BlockId: "blockid_123",
	}

	buffer, err := shared.ConvertObjectToJsonBuffer(originalObject)
	jsonString := string(buffer.Bytes())

	log.Printf("json: %s\nerror: %v", jsonString, err)
	log.Println("finished")
}

func doResponseTesting() {
	// jsonString := `
	// {
	// 	"BlockInfos": [
	// 		{
	// 			"BlockId": "myfile_1",
	// 			"DataNodeList": [
	// 				"1.1.1.1",
	// 				"2.2.2.2"
	// 			]
	// 		},
	// 		{
	// 			"BlockId": "myfile_2",
	// 			"DataNodeList": [
	// 				"3.3.3.3",
	// 				"4.4.4.4",
	// 				"5.5.5.5"
	// 			]
	// 		}
	// 	],
	// 	"Error": "an error occured"
	// }`

	jsonString := `
	{
		"Error": "an error occured"
	}`

	res := http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(jsonString)),
	}

	inflatedObject := shared.StoreBlockResponse{}
	err := shared.ObjectFromResponse(&res, &inflatedObject)

	log.Printf("object: %v\nerror: %v", inflatedObject, err)
	log.Println("finished")
}
