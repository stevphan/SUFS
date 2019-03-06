package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"shared"
	"strings"
)

var useLocalFile = false

const (
	// system constants
	blockSize int = 64 * 1024 * 1024 // 64MB

	// actions
	actionCreateFile string = "create-file"
	actionGetFile    string = "get-file"
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

func main() {
	doRequestTesting()
	return

	log.SetPrefix("- ")
	log.SetFlags(LogFlagFilenameAndLine)

	normalArgs, options := parseOsArgs()
	shared.Verbose = contains(options, "v")
	useLocalFile = contains(options, "use-local-file")

	if len(normalArgs) == 0 {
		log.Fatalf("Must supply an action of '%s' or '%s'\n", actionCreateFile, actionGetFile)
	}

	// TODO: ensure that tmp/create & tmp/get directories exist

	userAction := normalArgs[0]
	switch userAction {
	case actionCreateFile:
		shared.VerbosePrintln("User wants to create a file")
		createFile(normalArgs[1:])
	case actionGetFile:
		shared.VerbosePrintln("User wants to get a file")
		getFile(normalArgs[1:])
	default:
		log.Fatalf("Incorrect command. Must supply an action of '%s' or '%s'\n", actionCreateFile, actionGetFile)
	}
}

func parseOsArgs() (normalArgs []string, options []string) {
	normalArgs = []string{}
	options = []string{}

	for _, rawArg := range os.Args[1:] {
		isOption := strings.HasPrefix(rawArg, "-")
		if isOption {
			options = append(options, rawArg[1:])
		} else {
			normalArgs = append(normalArgs, rawArg)
		}
	}

	return
}
