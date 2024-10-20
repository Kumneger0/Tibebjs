package runtime

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	fileSystem "github.com/kumneger0/tibebjs/pkg/fs"
	console "github.com/kumneger0/tibebjs/pkg/globals"
	"github.com/kumneger0/tibebjs/pkg/modules"
	v8 "rogchap.com/v8go"
)

type Runtime struct {
	Isolate *v8.Isolate
	Context *v8.Context
}

func TransformScript(entryFilePath string) (string, error) {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{entryFilePath},
		Bundle:      true,
		Write:       false,
		Format:      api.FormatIIFE,
		GlobalName:  "global",
		LogLevel:    api.LogLevelInfo,
		Platform:    api.PlatformNode,
		Target:      api.ESNext,
			Plugins: []api.Plugin{
			{
				Name: "inject-dirname-filename",
				Setup: func(build api.PluginBuild) {
					build.OnLoad(api.OnLoadOptions{Filter: ".*"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
						fileContent, err := ioutil.ReadFile(args.Path)
						if err != nil {
							return api.OnLoadResult{}, err
						}
						dirPath := filepath.Dir(args.Path)
						fileName, err := filepath.Abs(args.Path)
						if err != nil {
							return api.OnLoadResult{}, err
						}

						injectedContent := fmt.Sprintf(`
							const __dirname = %q;
							const __filename = %q;
							%s
						`, dirPath, fileName, string(fileContent))

						return api.OnLoadResult{
							Contents:   &injectedContent,
							ResolveDir: dirPath,
						}, nil
					})
				},
			},
		},
	})

	if len(result.Errors) > 0 {
		return "", fmt.Errorf("error during esbuild transformation: %v", result.Errors)
	}
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
