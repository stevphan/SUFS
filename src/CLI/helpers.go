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

func sendRequestToNameNode(nameNodeAddr, path, method string, apiRequest interface{}, apiResponse interface{}) {
	buffer, err := shared.ConvertObjectToJsonBuffer(apiRequest)
	shared.CheckErrorAndFatal("Error while communicating with the name node", err)

	url := "http://" + nameNodeAddr + path

	client := http.Client{Timeout: nameNodeTimeout}
	req, err := http.NewRequest(method, url, buffer)
	shared.CheckErrorAndFatal("Error creating request to name node", err)

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	shared.CheckErrorAndFatal("Error while communicating with the name node", err)

	err = shared.ObjectFromResponse(res, apiResponse)
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
