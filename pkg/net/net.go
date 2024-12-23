package net

import (
	"fmt"

	"github.com/kumneger0/tibebjs/pkg/eventloop"
	v8 "rogchap.com/v8go"
)

type TibebRequest struct {
	Url string
}

type TibebResponse struct {
	Status     int
	Headers    map[string]string
	StatusText string
}

func Request(info *v8.FunctionCallbackInfo) *v8.Value {
	return v8.Undefined(info.Context().Isolate())
}

func Response(info *v8.FunctionCallbackInfo) *v8.Value {
	body := info.Args()[0]

	if body.IsNullOrUndefined() {
		panic("body is required")
	}
	options, err := info.Args()[1].AsObject()
	if err != nil {
		fmt.Println(err.Error())
	}

	status, err := options.Get("status")
	if err != nil {
		status, err = v8.NewValue(info.Context().Isolate(), 200)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	headers, err := options.Get("headers")
	if err != nil {
		fmt.Println(err.Error())
	}
	headerObj, err := headers.AsObject()
	if err != nil {
		fmt.Println(err.Error())
	}
	content, _ := headerObj.Value.MarshalJSON()
	v8obj := v8.NewObjectTemplate(info.Context().Isolate())
	v8obj.Set("body", body)
	v8obj.Set("status", status)
	v8obj.Set("statusText", "OK")
	v8obj.Set("headers", string(content))
	valueTobeReturned, _ := v8obj.NewInstance(info.Context())
	return valueTobeReturned.Value
}

var NetFuncs = []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
}{
	{
		Name: "response",
		Fn:   Response,
	},
	{
		Name: "request",
		Fn:   Request,
	},
}

func Serve(info *v8.FunctionCallbackInfo) *v8.Value {
	v8func, err := info.Args()[0].AsFunction()
	if err != nil {
		panic(err.Error())
	}
	eventloop.NetworkTaskQueue = append(eventloop.NetworkTaskQueue, eventloop.NetworkTask{
		Callback: v8func,
		Context:  info.Context(),
	})
	eventloop.Serve(info)
	return v8.Undefined(info.Context().Isolate())
}

var NetObj = []struct{
    Name string
    Fn func(*v8.FunctionCallbackInfo) *v8.Value
}{
    {
        Name: "serve",
        Fn:   Serve,
    },
}

func GetNetObjects() []struct{
    Name string
    Fn func(*v8.FunctionCallbackInfo) *v8.Value
}{
    return NetObj
}