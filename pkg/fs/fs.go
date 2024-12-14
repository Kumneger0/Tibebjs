package fs

import (
	"fmt"
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

	value, err := v8.NewValue(iso, string(fileContent))
	if err != nil {
		errMessage, _ := v8.NewValue(iso, err.Error())
		return errMessage
	}

	return value
}


func WriteFile(info *v8.FunctionCallbackInfo) *v8.Value {
	iso := info.Context().Isolate()
	path := info.Args()[0].String()

	var content []byte

	if info.Args()[1].IsString() {
		content = []byte(info.Args()[1].String())
	} else if info.Args()[1].IsObject() {
		obj := info.Args()[1].Object()

		if obj.IsTypedArray() {
			length, _ := obj.Get("length")
			size := length.Int32() 

			content = make([]byte, size)

			for i := int32(0); i < size; i++ {
				indexStr := fmt.Sprintf("%d", i) 
				val, _ := obj.Get(indexStr)      
				content[i] = byte(val.Int32())  
			}
		}
	}

	err := os.WriteFile(path, content, 0644)
	if err != nil {
		errMessage, _ := v8.NewValue(iso, err.Error())
		return errMessage
	}
	value, err := v8.NewValue(iso, "File written successfully")
	if err != nil {
		errMessage, _ := v8.NewValue(iso, err.Error())
		return errMessage
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
