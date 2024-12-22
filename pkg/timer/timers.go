package timer

import (
	eventloop "github.com/kumneger0/tibebjs/pkg/eventloop"
	v8 "rogchap.com/v8go"
)

func setTimeout(info *v8.FunctionCallbackInfo) *v8.Value {
	callback, _ := info.Args()[0].AsFunction()
	eventloop.TimerTaskQueue = append(eventloop.TimerTaskQueue, eventloop.TimerTask{
		Callback: callback,
		Context:  info.Context(),
	})
	eventloop.Schedule(info, false)
  	return v8.Undefined(info.Context().Isolate())
}

func clearTImeout(info *v8.FunctionCallbackInfo) *v8.Value {
    //TODO: implement clearTimeout
	return v8.Undefined(info.Context().Isolate())
}

func setInterval(info *v8.FunctionCallbackInfo) *v8.Value {
	callback, _ := info.Args()[0].AsFunction()
	eventloop.TimerTaskQueue = append(eventloop.TimerTaskQueue, eventloop.TimerTask{
		Callback: callback,
		Context:  info.Context(),
	})
	eventloop.Schedule(info, true)
	return v8.Undefined(info.Context().Isolate())
}

func clearInterval(info *v8.FunctionCallbackInfo) *v8.Value {
	//TODO: implement clearInterval
	return v8.Undefined(info.Context().Isolate())
}


var TimerFunctions = []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
}{
	{"setTimeout", setTimeout},
	{"clearTimeout", clearTImeout},
	{"setInterval", setInterval},
	{"clearInterval", clearInterval},
}

func GetTimerObjects() []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
} {
	return TimerFunctions
}
