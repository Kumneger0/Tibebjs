package timer

import (
	"fmt"

	eventloop "github.com/kumneger0/tibebjs/pkg/eventloop"
	v8 "rogchap.com/v8go"
)

const BASE_TIMEOUT_ID = 100
const BASE_INTERVAL_ID = 1000

func setTimeout(info *v8.FunctionCallbackInfo) *v8.Value {
	callback, _ := info.Args()[0].AsFunction()
	id := BASE_TIMEOUT_ID + len(eventloop.TimerTaskQueue)
	eventloop.Mu.Lock()
	defer eventloop.Mu.Unlock()
	eventloop.TimerTaskQueue = append(eventloop.TimerTaskQueue, eventloop.TimerTask{
		Callback: callback,
		Context:  info.Context(),
		Id:       id,
	})
	eventloop.Schedule(info, false, id)
	idToReturnToJs, err := v8.NewValue(info.Context().Isolate(), float64(id))
	if err != nil {
		fmt.Println("error creating id to return to js", err.Error())
	}
	return idToReturnToJs
}

func clearTimeout(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("clearing timeout")
	id := int(info.Args()[0].Int32())
	task, err := eventloop.GetTask(id)
	if err != nil {
		return v8.Undefined(info.Context().Isolate())
	}
	task.Clear()
	return v8.Undefined(info.Context().Isolate())
}

func setInterval(info *v8.FunctionCallbackInfo) *v8.Value {
	callback, _ := info.Args()[0].AsFunction()
	id := BASE_INTERVAL_ID + len(eventloop.TimerTaskQueue)
	eventloop.Mu.Lock()
	defer eventloop.Mu.Unlock()
	eventloop.TimerTaskQueue = append(eventloop.TimerTaskQueue, eventloop.TimerTask{
		Callback: callback,
		Context:  info.Context(),
		Id:       id,
	})
	eventloop.Schedule(info, true, id)
	idToReturnToJs, err := v8.NewValue(info.Context().Isolate(), float64(id))
	if err != nil {
		fmt.Println("error creating id to return to js", err.Error())
	}
	return idToReturnToJs
}

func clearInterval(info *v8.FunctionCallbackInfo) *v8.Value {
	fmt.Println("clearing interval")
	id := int(info.Args()[0].Int32())
	task, err := eventloop.GetTask(id)
	if err != nil {
		return v8.Undefined(info.Context().Isolate())
	}
	task.Clear()
	return v8.Undefined(info.Context().Isolate())
}

var TimerFunctions = []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
}{
	{"setTimeout", setTimeout},
	{"clearTimeout", clearTimeout},
	{"setInterval", setInterval},
	{"clearInterval", clearInterval},
}

func GetTimerObjects() []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
} {
	return TimerFunctions
}
