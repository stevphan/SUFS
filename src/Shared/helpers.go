package Shared

import (
	"bytes"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"log"
	"fmt"
)

var Verbose = false

func VerbosePrintln(s string) {
	if Verbose {
		fmt.Println(s)
	}
}

func CheckErrorAndFatal(description string, err error) {
	if err != nil {
		log.Fatal(description+":", err)
	}
}

func ConvertObjectToJsonBuffer(object interface{}) (*bytes.Buffer, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(data)

	return buffer, nil
}

func ObjectFromResponse(res *http.Response, object interface{}) error {
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
