package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var verbose = false
var useLocalFile = false

const (
	// actions
	actionCreateFile string = "create-file"
	actionGetFile    string = "get-file"

	// directories for temporary files
	dirTempCreateFiles string = "/Users/Rivukis/Desktop/tmp/create/"
	dirTempGetFiles    string = "/Users/Rivukis/Desktop/tmp/get/"
)

func doRequestTesting() {
	// originalObject := createFileNameNodeRequest{
	// 	FileName: "myfilename",
	// 	Size: "12345",
	// }

	dnList := []string{"1.2.3.4", "5.6.7.8", "10.0.0.7"}

	originalObject := storeBlockRequest{
		Block:   "actualblock",
		DnList:  dnList,
		BlockId: "blockid_123",
	}

	buffer, err := convertObjectToJsonBuffer(originalObject)
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

	inflatedObject := storeBlockResponse{}
	err := objectFromResponse(&res, &inflatedObject)

	log.Printf("object: %v\nerror: %v", inflatedObject, err)
	log.Println("finished")
}

func main() {
	log.SetPrefix("- ")
	log.SetFlags(LogFlagFilenameAndLine)

	normalArgs, options := parseOsArgs()
	verbose = contains(options, "v")
	useLocalFile = contains(options, "use-local-file")

	if len(normalArgs) == 0 {
		log.Fatalf("Must supply an action of '%s' or '%s'\n", actionCreateFile, actionGetFile)
	}

	// TODO: ensure that tmp/create & tmp/get directories exist

	userAction := normalArgs[0]
	switch userAction {
	case actionCreateFile:
		verbosePrintln("User wants to create a file")
		createFile(normalArgs[1:])
	case actionGetFile:
		verbosePrintln("User wants to get a file")
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
