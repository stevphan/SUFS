package main

import (
	"net/http"
	"shared"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func stringsMap(strs []string, mapper func(string) string) []string {
	result := []string{}
	for _, str := range strs {
		result = append(result, mapper(str))
	}

	return result
}

func sendRequestToNameNode(nameNodeAddr, path string, request interface{}, response interface{}) {
	buffer, err := shared.ConvertObjectToJsonBuffer(request)
	shared.CheckErrorAndFatal("Error while communicating with the name node:", err)

	url := "http://" + nameNodeAddr + "/" + path

	client := http.Client{Timeout: nameNodeTimeout}
	res, err := client.Post(url, "application/json", buffer)
	shared.CheckErrorAndFatal("Error while communicating with the name node:", err)

	err = shared.ObjectFromResponse(res, response)
	shared.CheckErrorAndFatal("Unable to parse response", err)

	return
}

type LogFlag int

const (
	LogFlagDate            = 1
	LogFlagTime            = 2
	LogFlagTimeDecimal     = 4
	LogFlagFilePathAndLine = 8
	LogFlagFilenameAndLine = 16
)
