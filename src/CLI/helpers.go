package main

import (
	"fmt"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func verbosePrintln(s string) {
	if verbose {
		fmt.Println(s)
	}
}

type LogFlag int

const (
	LogFlagDate = 1
	LogFlagTime = 2
	LogFlagTimeDecimal = 4
	LogFlagFilePathAndLine = 8
	LogFlagFilenameAndLine = 16
)
