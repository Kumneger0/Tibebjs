package console

import (
	"fmt"
	"strings"

	v8 "rogchap.com/v8go"
)

func Log(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println(strings.Join(args, " "))
	return nil
}

func Info(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("INFO:", strings.Join(args, " "))
	return nil
}

func Warn(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.warn is not implemented yet.")
	return nil
}

func Error(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.error is not implemented yet.")
	return nil
}

func Debug(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.debug is not implemented yet.")
	return nil
}

func Assert(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.assert is not implemented yet.")
	return nil
}

func Clear(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.clear is not implemented yet.")
	return nil
}

func Count(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.count is not implemented yet.")
	return nil
}

func CountReset(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.countReset is not implemented yet.")
	return nil
}

func Group(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.group is not implemented yet.")
	return nil
}

func GroupEnd(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.groupEnd is not implemented yet.")
	return nil
}

func Table(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.table is not implemented yet.")
	return nil
}

func Time(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.time is not implemented yet.")
	return nil
}

func TimeEnd(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.timeEnd is not implemented yet.")
	return nil
}

func Trace(info *v8.FunctionCallbackInfo) *v8.Value {
	args := make([]string, len(info.Args()))
	for i, arg := range info.Args() {
		args[i] = arg.String()
	}
	fmt.Println("console.trace is not implemented yet.")
	return nil
}

var consoleFunctions = []struct {
	name string
	fn   func(*v8.FunctionCallbackInfo) *v8.Value
}{
	{"log", Log},
	{"error", Error},
	{"warn", Warn},
	{"info", Info},
	{"debug", Debug},
	{"group", Group},
	{"groupEnd", GroupEnd},
	{"table", Table},
	{"time", Time},
	{"timeEnd", TimeEnd},
	{"trace", Trace},
	{"assert", Assert},
	{"clear", Clear},
	{"count", Count},
	{"countReset", CountReset},
}

func CreateConsoleObject(iso *v8.Isolate) *v8.ObjectTemplate {
	console := v8.NewObjectTemplate(iso)

	for _, fn := range consoleFunctions {
		console.Set(fn.name, v8.NewFunctionTemplate(iso, fn.fn))
	}
	return console
}
