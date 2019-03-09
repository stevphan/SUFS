package main

import (
	"log"
	"os"
	"shared"
	"strings"
	"time"
)

const (
	// system constants
	blockSize       int = 64 * 1024 * 1024 // 64MB
	nameNodeTimeout     = time.Duration(5 * time.Second)

	// actions
	actionCreateFile    string = "create-file"
	actionGetFile       string = "get-file"
	actionListDataNodes string = "list-data-nodes"

	// S3 constants
	// tempS3DownloadFileName string = "temp_s3_download_file_name"
	awsRegion string = "us-west-2"
)

var (
	awsAccessId          string
	awsSecretAccessToken string
)

func main() {
	log.SetPrefix("- ")
	log.SetFlags(LogFlagFilenameAndLine)

	normalArgs, options := parseArgs()
	shared.Verbose = contains(options, "v")

	parseEnvironmentVariables()

	if len(normalArgs) == 0 {
		log.Fatalf("Must supply an action of '%s', '%s', or '%s'\n", actionCreateFile, actionGetFile, actionListDataNodes)
	}

	userAction := normalArgs[0]
	switch userAction {
	case actionCreateFile:
		shared.VerbosePrintln("User wants to create a file")
		createFile(normalArgs[1:])
	case actionGetFile:
		shared.VerbosePrintln("User wants to get a file")
		getFile(normalArgs[1:], false)
	case actionListDataNodes:
		shared.VerbosePrintln("User wants to get Data Node info of a file")
		getFile(normalArgs[1:], true)
	default:
		log.Fatalf("Incorrect command. Must supply an action '%s', '%s' or '%s'\n", actionCreateFile, actionGetFile, actionListDataNodes)
	}
}

func parseArgs() (normalArgs []string, options []string) {
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

func parseEnvironmentVariables() {
	awsAccessId = os.Getenv("AWS_ACCESS_ID")
	awsSecretAccessToken = os.Getenv("AWS_SECRET_ACCESS_TOKEN")

	if len(awsAccessId) == 0 {
		log.Fatal("Environment variable AWS_ACCESS_ID not set")
	}

	if len(awsSecretAccessToken) == 0 {
		log.Fatal("Environment variable AWS_SECRET_ACCESS_TOKEN not set")
	}
}
