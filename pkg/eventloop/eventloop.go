package eventloop

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	utils "github.com/kumneger0/tibebjs/pkg/utils"
	v8 "rogchap.com/v8go"
)

type TimerTask struct {
	Callback *v8.Function
	Context  *v8.Context
	Id       int
	Cleared  bool
}

type NetworkTask struct {
	Callback *v8.Function
	Context  *v8.Context
	FuncArg  *v8.Value
}

var (
	Mu             sync.RWMutex
	TimerTaskQueue []TimerTask
)

var TimerTaskChannel chan TimerTask

var NetworkTaskQueue []NetworkTask
var NetworkTaskChannel chan NetworkTask
var NetworkTaskResponseChannel chan *v8.Value

func init() {
	TimerTaskChannel = make(chan TimerTask, 100)
	NetworkTaskChannel = make(chan NetworkTask, 100)
	NetworkTaskResponseChannel = make(chan *v8.Value, 100)
	go func() {
		for {
			if len(TimerTaskQueue) == 0 {
				return
			}
		}
	}()
}

func (t *TimerTask) Clear() {
	Mu.Lock()
	defer Mu.Unlock()
	t.Cleared = true

	for i, task := range TimerTaskQueue {
		if task.Id == t.Id {
			TimerTaskQueue = append(TimerTaskQueue[:i], TimerTaskQueue[i+1:]...)
			break
		}
	}
}

func (t *TimerTask) IsCleared() bool {
	Mu.RLock()
	defer Mu.RUnlock()
	return t.Cleared
}

func (t *TimerTask) Add() bool {
	Mu.Lock()
	defer Mu.Unlock()
	TimerTaskQueue = append(TimerTaskQueue, *t)
	return true
}

func GetTask(id int) (*TimerTask, error) {
	var TimerTask *TimerTask
	for _, task := range TimerTaskQueue {
		if task.Id == id {
			TimerTask = &task
		}
	}
	if TimerTask == nil {
		return nil, errors.New("task not found")
	}
	return TimerTask, nil
}

func Schedule(info *v8.FunctionCallbackInfo, interval bool, id int) {
	var _ = make([]string, len(info.Args()))
	callback := info.Args()[0]
	if !callback.IsFunction() {
		fmt.Println("the first argument must be function")
		os.Exit(1)
	}
	_, err := callback.AsFunction()
	if err != nil {
		fmt.Println("the first argument must be function")
		os.Exit(1)
	}
	delay := info.Args()[1].Int32()
	if !interval {
		go func() {
			time.Sleep(time.Duration(delay) * time.Millisecond)
			task, error := GetTask(id)
			if error != nil || task.IsCleared() {
				return
			}
			TimerTaskChannel <- *task
			TimerTaskQueue = TimerTaskQueue[:len(TimerTaskQueue)-1]
		}()
	} else {
		go func() {
			for {
				task, error := GetTask(id)
				if error != nil || task.IsCleared() {
					return
				}
				time.Sleep(time.Duration(delay) * time.Millisecond)
				TimerTaskChannel <- *task
			}
		}()
	}
}

func Serve(info *v8.FunctionCallbackInfo) {
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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		request := utils.MakeJSRequestObj(r, info)

		NetworkTaskChannel <- NetworkTask{
			Callback: HandleFunc,
			Context:  info.Context(),
			FuncArg:  request,
		}

		NetworkTaskQueue = append(NetworkTaskQueue, NetworkTask{
			Callback: HandleFunc,
			Context:  info.Context(),
			FuncArg:  request,
		})

		res := <-NetworkTaskResponseChannel
		extractedResponse := utils.ExteactResponse(res)

		body := extractedResponse.Body
		status := extractedResponse.Status
		statusText := extractedResponse.StatusText
		var headers map[string]interface{}
		json.Unmarshal([]byte(extractedResponse.Headers), &headers)
		for key, value := range headers {
			w.Header().Set(key, value.(string))
		}
		w.WriteHeader(status)
		w.Header().Set("statusText", statusText)
		_, err := w.Write([]byte(body))

		if err != nil {
			fmt.Println("Error writing response:", err)
		}
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
}
