package fs

import (
	"fmt"

	eventloop "github.com/kumneger0/tibebjs/pkg/eventloop"

	v8 "rogchap.com/v8go"
)

func ReadFile(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("reading file")
	promise := eventloop.ReadFile(info)
	return promise.Value
}

func WriteFile(info *v8.FunctionCallbackInfo) *v8.Value {
	return eventloop.WriteFile(info).Value
}

func RenameFile(info *v8.FunctionCallbackInfo) *v8.Value {
	return eventloop.RenameFile(info).Value
}

func serve(info *v8.FunctionCallbackInfo) *v8.Value {
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

func RMFile(info *v8.FunctionCallbackInfo) *v8.Value {
	return eventloop.RmFile(info).Value
}



var fsFunctions = []struct {
	name string
	fn   func(*v8.FunctionCallbackInfo) *v8.Value
}{
	{"readFile", ReadFile},
	{"writeFile", WriteFile},
	{"rmFile", RMFile},
	{"serve", serve},
	{"renameFile", RenameFile},
}

func CreateFsObject(iso *v8.Isolate) *v8.ObjectTemplate {
	fs := v8.NewObjectTemplate(iso)
	for _, fn := range fsFunctions {
		fs.Set(fn.name, v8.NewFunctionTemplate(iso, fn.fn))
	}
	return fs
}
