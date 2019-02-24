package main

import (
	"log"
	"os"
	"strings"
)

var verbose = false

const (
	actionCreateFile string = "create-file"
	actionGetFile    string = "get-file"
)

func main() {
	log.SetPrefix("- ")
	log.SetFlags(LogFlagFilenameAndLine)

	normalArgs, options := parseOsArgs()
	verbose = contains(options, "v")

	if len(normalArgs) == 0 {
		log.Fatalf("Must supply an action of '%s' or '%s'\n", actionCreateFile, actionGetFile)
	}

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
