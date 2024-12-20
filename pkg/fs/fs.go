package fs

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

func serve(info *v8.FunctionCallbackInfo) *v8.Value {
	if len(info.Args()) < 2 {
		panic("serve requires at least 2 arguments")
	}
	if !info.Args()[0].IsFunction() {
		panic("The first argument must be a function")
	}

	HandleFunc, err := info.Args()[0].AsFunction()
	if err != nil {
		panic(err.Error())
	}

	port := info.Args()[1].Int32()
	mux := http.NewServeMux()

	requestResChannel := make(chan struct {
		r *http.Request
		w http.ResponseWriter
	}, 100)

	response := make(chan string, 100)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestResChannel <- struct {
			r *http.Request
			w http.ResponseWriter
		}{r, w}

		res := <-response
		result, err := w.Write([]byte(res))

		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		fmt.Printf("Wrote %d bytes\n", result)

	})
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	serverReady := make(chan struct{})

	go func() {
		fmt.Printf("Server started on port %d\n", port)
		close(serverReady)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err.Error())
		}
	}()

	<-serverReady

	for {
		select {
		case request := <-requestResChannel:
			w := request.w
			r := request.r

			fmt.Printf("Request received: %s\n", r.URL.Path)
			value, err := HandleFunc.Call(v8.Undefined(info.Context().Isolate()))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error processing request"))
				continue
			}
			response <- value.String()

		case <-shutdownChan:
			fmt.Println("Shutting down server...")
			if err := server.Shutdown(context.Background()); err != nil {
				fmt.Printf("Server shutdown error: %v\n", err)
			}
		}
	}

}

var fsFunctions = []struct {
	name string
	fn   func(*v8.FunctionCallbackInfo) *v8.Value
}{
	{"readFile", ReadFile},
	{"writeFile", WriteFile},
	{"serve", serve},
}

func CreateFsObject(iso *v8.Isolate) *v8.ObjectTemplate {
	fs := v8.NewObjectTemplate(iso)
	for _, fn := range fsFunctions {
		fs.Set(fn.name, v8.NewFunctionTemplate(iso, fn.fn))
	}
	return fs
}
