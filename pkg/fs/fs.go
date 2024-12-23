package fs

import (
	eventloop "github.com/kumneger0/tibebjs/pkg/eventloop"
	v8 "rogchap.com/v8go"
)

func ReadFile(info *v8.FunctionCallbackInfo) *v8.Value {
	return eventloop.ReadFile(info).Value
}

func WriteFile(info *v8.FunctionCallbackInfo) *v8.Value {
	return eventloop.WriteFile(info).Value
}

func RenameFile(info *v8.FunctionCallbackInfo) *v8.Value {
	return eventloop.RenameFile(info).Value
}

func RMFile(info *v8.FunctionCallbackInfo) *v8.Value {
	return eventloop.RmFile(info).Value
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

func GetFsObjects() []struct{
    Name string
    Fn func(*v8.FunctionCallbackInfo) *v8.Value
} {
	return fsFunctions
}