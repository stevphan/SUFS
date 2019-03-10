package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var Verbose = false

func VerbosePrintln(s string) {
	if Verbose {
		fmt.Println(s)
	}
}

func CheckErrorAndFatal(description string, err error) {
	if err != nil {
		log.Fatalln(description+":", err)
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

func ServeCall(pattern string, handlers map[string]func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handler, present := handlers[r.Method]
		if present {
			handler(w, r)
		}
	})
}
