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

func GoValueToV8(isolate *v8.Isolate, value interface{}, ctx *v8.Context) (*v8.Value, error) {
	switch v := value.(type) {
	case string:
		return v8.NewValue(isolate, v)
	case float64:
		return v8.NewValue(isolate, v)
	case bool:
		return v8.NewValue(isolate, v)
	case []interface{}:
		arrayTmpl := v8.NewObjectTemplate(isolate)
		array, err := arrayTmpl.NewInstance(ctx)
		if err != nil {
			return nil, err
		}
		// Set array length
		array.Set("length", len(v))
		for i, item := range v {
			itemValue, err := GoValueToV8(isolate, item, ctx)
			if err != nil {
				return nil, err
			}
			array.Set(fmt.Sprint(i), itemValue)
		}
		return array.Value, nil
	case map[string]interface{}:
		obj := v8.NewObjectTemplate(isolate)
		instance, err := obj.NewInstance(ctx)
		if err != nil {
			return nil, err
		}
		for key, item := range v {
			itemValue, err := GoValueToV8(isolate, item, ctx)
			if err != nil {
				return nil, err
			}
			instance.Set(key, itemValue)
		}
		return instance.Value, nil
	case nil:
		return v8.Null(isolate), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}

func Text(info *v8.FunctionCallbackInfo, response *http.Response) *v8.Value {
	textPromiseResolver, err := v8.NewPromiseResolver(info.Context())
	if err != nil {
		panic(err.Error())
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		errValue, _ := v8.NewValue(info.Context().Isolate(), fmt.Sprintf("Failed to read body: %s", err.Error()))
		textPromiseResolver.Reject(errValue)
		return textPromiseResolver.GetPromise().Value
	}

	text, err := v8.NewValue(info.Context().Isolate(), string(body))
	if err != nil {
		panic(err.Error())
	}

	textPromiseResolver.Resolve(text)
	return textPromiseResolver.GetPromise().Value
}



 func Json(info *v8.FunctionCallbackInfo, response *http.Response) *v8.Value {
		jsonPromiseResolver, _ := v8.NewPromiseResolver(info.Context())

				fmt.Println("json() method called")
				body, err := io.ReadAll(response.Body)
				if err != nil {
					fmt.Printf("Error reading body: %v\n", err)
					errValue, _ := v8.NewValue(info.Context().Isolate(), fmt.Sprintf("Failed to read body: %s", err.Error()))
					jsonPromiseResolver.Reject(errValue)
					return jsonPromiseResolver.GetPromise().Value
				}

				var result interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					fmt.Printf("JSON parse error: %v\n", err)
					errValue, _ := v8.NewValue(info.Context().Isolate(), fmt.Sprintf("Failed to parse JSON: %s", err.Error()))
					jsonPromiseResolver.Reject(errValue)
					return jsonPromiseResolver.GetPromise().Value
				}

				fmt.Printf("Parsed JSON result: %+v\n", result)

				jsonValue, err := GoValueToV8(info.Context().Isolate(), result, info.Context())
				if err != nil {
					fmt.Printf("JSON conversion error: %v\n", err)
					errValue, _ := v8.NewValue(info.Context().Isolate(), fmt.Sprintf("Failed to convert JSON: %s", err.Error()))
					jsonPromiseResolver.Reject(errValue)
					return jsonPromiseResolver.GetPromise().Value
				}

				fmt.Printf("Successfully converted to V8 value\n")
				jsonPromiseResolver.Resolve(jsonValue)
				return jsonPromiseResolver.GetPromise().Value
			}




