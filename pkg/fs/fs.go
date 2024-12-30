package fs

import (
	"math/rand"

	eventloop "github.com/kumneger0/tibebjs/pkg/eventloop"
	v8 "rogchap.com/v8go"
)

func ReadFile(info *v8.FunctionCallbackInfo) *v8.Value {
	iotask := eventloop.IOTask{Callback: &v8.Function{}, Context: info.Context(), Id: rand.Int()}
	eventloop.IoTask = append(eventloop.IoTask, iotask)
	return eventloop.ReadFile(info, &iotask).Value
}

func WriteFile(info *v8.FunctionCallbackInfo) *v8.Value {
	iotask := eventloop.IOTask{Callback: &v8.Function{}, Context: info.Context(), Id: rand.Int()}
	eventloop.IoTask = append(eventloop.IoTask, iotask)
	return eventloop.WriteFile(info, &iotask).Value
}

func RenameFile(info *v8.FunctionCallbackInfo) *v8.Value {
	iotask := eventloop.IOTask{Callback: &v8.Function{}, Context: info.Context(), Id: rand.Int()}
	eventloop.IoTask = append(eventloop.IoTask, iotask)
	return eventloop.RenameFile(info, &iotask).Value
}

func RMFile(info *v8.FunctionCallbackInfo) *v8.Value {
	iotask := eventloop.IOTask{Callback: &v8.Function{}, Context: info.Context(), Id: rand.Int()}
	eventloop.IoTask = append(eventloop.IoTask, iotask)
	return eventloop.RmFile(info, &iotask).Value
}

var fsFunctions = []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
}{
	{"readFile", ReadFile},
	{"writeFile", WriteFile},
	{"rmFile", RMFile},
	{"renameFile", RenameFile},
}

func GetFsObjects() []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
} {
	return fsFunctions
}
