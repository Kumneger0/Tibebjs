package modules

import (
	"fmt"
	"os"
	"path/filepath"

	v8 "rogchap.com/v8go"
)

func CreateImportFunction(iso *v8.Isolate, baseDir string) (*v8.Function, error) {
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		args := info.Args()
		if len(args) < 1 {
			return throwException(iso, "import() requires at least one argument")
		}
		modulePath := args[0].String()
		fullPath := resolveModulePath(baseDir, modulePath)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			return throwException(iso, fmt.Sprintf("Error reading module file: %v", err))
		}

		moduleCtx := v8.NewContext(iso)
		defer moduleCtx.Close()

		importMetaObj := v8.NewObjectTemplate(iso)
		importMetaObj.Set("url", "file://"+fullPath)

		global := moduleCtx.Global()
		global.Set("import", info.This())
		importMetaInstance, _ := importMetaObj.NewInstance(moduleCtx)
		global.Set("import.meta", importMetaInstance)

		wrappedCode := fmt.Sprintf(`
			async function __module__() {
				const module = { exports: {} };
				const exports = module.exports;
				%s
				return module.exports;
			}
			__module__();
		`, string(content))

		result, err := moduleCtx.RunScript(wrappedCode, fullPath)
		if err != nil {
			return throwException(iso, fmt.Sprintf("Error executing module: %v", err))
		}

		promise, _ := v8.NewPromiseResolver(info.Context())
		promise.Resolve(result)

		return promise.GetPromise().Value
	}).GetFunction(v8.NewContext(iso)), nil
}

func resolveModulePath(baseDir, modulePath string) string {
	fullPath := filepath.Join(baseDir, modulePath)

	if !filepath.IsAbs(fullPath) {
		fullPath = filepath.Join(baseDir, fullPath)
	}

	// Add .js extension if not present
	if filepath.Ext(fullPath) == "" {
		fullPath += ".js"
	}

	return fullPath
}

func throwException(iso *v8.Isolate, message string) *v8.Value {
	exception, _ := v8.NewValue(iso, message)
	iso.ThrowException(exception)
	return nil
}
