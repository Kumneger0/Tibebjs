package main

import (
	"fmt"
	"os"
	"path/filepath"

	mainEventloop "github.com/kumneger0/tibebjs/pkg"
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

	_, err = rt.ExecuteScript(scriptPath)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	mainEventloop.Wg.Wait()
	for {
		select {
		case <-mainEventloop.Channel:
			{
				for _, task := range mainEventloop.Tasks {
					task.Callback.Call(v8.Undefined(rt.Isolate))
				}
				mainEventloop.Tasks = make([]mainEventloop.Task, 0)
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
