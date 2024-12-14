package eventloop

import (
	"fmt"
	"os"
	"time"

	v8 "rogchap.com/v8go"
)

func setTimeout(info *v8.FunctionCallbackInfo) *v8.Value {
	var _ = make([]string, len(info.Args()))

	callback := info.Args()[0]

	if !callback.IsFunction() {
		fmt.Println("the first argument must be function")
		os.Exit(1)
	}
	v8Function, err := callback.AsFunction()
	if err != nil {
		fmt.Println("the first argument must be function")
		os.Exit(1)
	}

	delay := info.Args()[1].Int32()

	timer := time.After(time.Duration(delay) * time.Millisecond)

	fmt.Println("about to start timer")

	 <-timer
	v8Function.Call(v8.Undefined(info.Context().Isolate()))
 
	return v8.Undefined(info.Context().Isolate())
}

var TimerFunctions = []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
}{
	{"setTimeout", setTimeout},
}

func GetTimerObjects() []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
} {
	return TimerFunctions
}
