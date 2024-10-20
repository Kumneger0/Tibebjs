package fs

import (
	"os"

	v8 "rogchap.com/v8go"
)

func ReadFile(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	path := info.Args()[0].String()

	fileContent, err := os.ReadFile(path)
	if err != nil {
		errMessage, _ := v8.NewValue(iso, err.Error())
		return errMessage
	}
	value, err := v8.NewValue(iso, fileContent)
	if err != nil {
		err, _ := v8.NewValue(iso, err.Error())
		return err
	}

	return value
}

func WriteFile(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()

	path := info.Args()[0].String()
	content := info.Args()[1].String()
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		errMessage, _ := v8.NewValue(iso, err.Error())
		return errMessage
	}
	value, err := v8.NewValue(iso, "File written successfully")
	if err != nil {
		err, _ := v8.NewValue(iso, err.Error())
		return err
	}

	return value
}

var fsFunctions = []struct {
	name string
	fn   func(*v8.FunctionCallbackInfo) *v8.Value
}{
	{"readFile", ReadFile},
	{"writeFile", WriteFile},
}

func CreateFsObject(iso *v8.Isolate) *v8.ObjectTemplate {
	fs := v8.NewObjectTemplate(iso)

	for _, fn := range fsFunctions {
		fs.Set(fn.name, v8.NewFunctionTemplate(iso, fn.fn))
	}
	return fs
}
