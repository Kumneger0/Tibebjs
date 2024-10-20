package runtime

import (
	"fmt"
	"os"

	"github.com/evanw/esbuild/pkg/api"
	fileSystem "github.com/kumneger0/tibebjs/pkg/fs"
	console "github.com/kumneger0/tibebjs/pkg/globals"
	"github.com/kumneger0/tibebjs/pkg/modules"
	v8 "rogchap.com/v8go"
)

var ModuleCache = make(map[string]*v8.Value)

type Runtime struct {
	Isolate *v8.Isolate
	Context *v8.Context
}

func TransformScript(entryFilePath string) (string, error) {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{entryFilePath}, // Entry file (main.js)
		Bundle:      true,                    // Bundle everything together
		Write:       false,                   // Do not write to disk
		Format:      api.FormatIIFE,          // Use IIFE format (immediately invoked function)
		GlobalName:  "global",                // Name of the global object (optional)
		LogLevel:    api.LogLevelInfo,
		Platform:    api.PlatformNode, // Simulate Node.js environment
		Target:      api.ESNext,       // Target latest JS
	})

	// Check if esbuild encountered any errors
	if len(result.Errors) > 0 {
		return "", fmt.Errorf("error during esbuild transformation: %v", result.Errors)
	}

	// Extract the transformed/bundled code
	bundledCode := string(result.OutputFiles[0].Contents)

	return bundledCode, nil
}

func NewRuntime() (*Runtime, error) {
	iso := v8.NewIsolate()
	if iso == nil {
		return nil, fmt.Errorf("failed to create isolate")
	}

	ctx := v8.NewContext(iso)
	return &Runtime{Isolate: iso, Context: ctx}, nil
}

func (r *Runtime) SetupGlobals(scriptDir string) error {
	global := r.Context.Global()

	fileSystemApi, err := fileSystem.CreateFsObject(r.Isolate).NewInstance(r.Context)

	if err != nil {
		fmt.Println(err.Error())
	}
	global.Set("Tibeb", fileSystemApi)

	consoleObj, err := console.CreateConsoleObject(r.Isolate).NewInstance(r.Context)
	if err != nil {
		return fmt.Errorf("error creating console instance: %v", err)
	}
	err = global.Set("console", consoleObj)
	if err != nil {
		return fmt.Errorf("error setting console object: %v", err)
	}

	importFn, err := modules.CreateImportFunction(r.Isolate, scriptDir)
	if err != nil {
		return fmt.Errorf("error creating import function: %v", err)
	}
	err = global.Set("__import__", importFn)
	if err != nil {
		return fmt.Errorf("error setting import function: %v", err)
	}

	return nil
}

func (r *Runtime) ExecuteScript(scriptPath string) (*v8.Value, error) {
	_, err := os.ReadFile(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("error reading script file: %v", err)
	}

	transformedScript, err := TransformScript(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("error transforming script: %v", err)
	}

	wrappedScript := fmt.Sprintf(`
		(async () => {
			%s
		})();
	`, transformedScript)

	return r.Context.RunScript(wrappedScript, scriptPath)
}

func (r *Runtime) Dispose() {
	r.Context.Close()
	r.Isolate.Dispose()
}
