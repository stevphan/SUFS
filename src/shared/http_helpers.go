package shared

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	PathFile        string = "/file"        // GET, PUT
	PathBlock       string = "/block"       // GET, PUT
	PathBlockReport string = "/blockReport" // PUT
	PathHeartBeat   string = "/heartBeat"   // PUT
)

func ServeCall(pattern string, handlers map[string]func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		handler, present := handlers[r.Method]
		if present {
			handler(w, r)
		} else {
			log.Printf("Unable to find handler for %s %s\n", r.Method, pattern)
		}
	})
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
