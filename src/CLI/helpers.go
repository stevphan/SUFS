package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func verbosePrintln(s string) {
	if verbose {
		fmt.Println(s)
	}
}

func checkErrorAndFatal(description string, err error) {
	if err != nil {
		log.Fatal(description+":", err)
	}
}

func convertObjectToJsonBuffer(object interface{}) (*bytes.Buffer, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(data)

	return buffer, nil
}

func objectFromResponse(res *http.Response, object interface{}) error {
	defer res.Body.Close()

	data := []byte{}
	if res.Body != nil {

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		data = body
	}

	err := json.Unmarshal(data, object)
	return err
}

type LogFlag int

const (
	LogFlagDate            = 1
	LogFlagTime            = 2
	LogFlagTimeDecimal     = 4
	LogFlagFilePathAndLine = 8
	LogFlagFilenameAndLine = 16
)
