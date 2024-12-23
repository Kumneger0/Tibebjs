package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	eventloop "github.com/kumneger0/tibebjs/pkg/eventloop"
	runtime "github.com/kumneger0/tibebjs/pkg/runtime"
	v8 "rogchap.com/v8go"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	scriptPath, err := getScriptPath()
	if err != nil {
		return err
	}
	rt, err := runtime.NewRuntime()
	if err != nil {
		return err
	}

	defer rt.Dispose()
	scriptDir := filepath.Dir(scriptPath)
	if err := rt.SetupGlobals(scriptDir); err != nil {
		return err
	}

	result, err := rt.ExecuteScript(scriptPath)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	promise, err := result.AsPromise()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	result = promise.Result()
	if !result.IsNullOrUndefined() {
		jsError := v8.JSError{
			Message:    result.String(),
			StackTrace: scriptPath,
		}

		log.Fatalf("\n %s\n Stack Trace: %s\n",
			jsError.Message,
			jsError.StackTrace,
		)
	}

	for {
		select {
		case task := <-eventloop.TimerTaskChannel:
			{
				task.Callback.Call(v8.Undefined(rt.Isolate))
			}
		case networkTask := <-eventloop.NetworkTaskChannel:
			{
				value, err := networkTask.Callback.Call(v8.Undefined(rt.Isolate), networkTask.FuncArg)
				if err != nil {
					panic(err.Error())
				}
				eventloop.NetworkTaskResponseChannel <- value
			}
		}
	}
}

func getScriptPath() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("usage: go run main.go <path_to_script>")
	}
	return os.Args[1], nil
}
