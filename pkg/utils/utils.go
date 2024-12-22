package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	v8 "rogchap.com/v8go"
)

func Log(message string) {
	log.Println(message)
}

type ResponseType struct {
	Body       string `json:"body"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
	Headers    string `json:"headers"`
}

func ExteactResponse(res *v8.Value) ResponseType {
	response := res.Object()
	body, err := response.Get("body")
	if err != nil {
		fmt.Println(err.Error())
	}
	status, _ := response.Get("status")
	statusText, _ := response.Get("statusText")
	headers, _ := response.Get("headers")

	statusNum := status.Int32()

	if !status.IsNullOrUndefined() {
		statusNum = status.Int32()
	}

	return ResponseType{
		Body:       body.String(),
		Status:     int(statusNum),
		StatusText: statusText.String(),
		Headers:    headers.String(),
	}
}
func MakeJSRequestObj(r *http.Request, info *v8.FunctionCallbackInfo) *v8.Value {
	url := r.URL.String()
	method := r.Method
	headers := r.Header
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	headersmap := make(map[string][]interface{})
	for key, values := range headers {
		for _, value := range values {
			headersmap[key] = append(headersmap[key], value)
		}
	}

	headersjson, err := json.Marshal(headersmap)
	if err != nil {
		panic(err.Error())
	}

	requestObj := v8.NewObjectTemplate(info.Context().Isolate())
	requestObj.Set("url", url)
	requestObj.Set("method", method)
	requestObj.Set("headers", string(headersjson))
	requestObj.Set("body", string(body))

	requestIntanace, err := requestObj.NewInstance(info.Context())
	if err != nil {
		fmt.Println(err.Error())
	}
	return requestIntanace.Value
}
