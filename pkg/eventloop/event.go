package eventloop

import (
	"fmt"
	"sync"
	"time"

	v8 "rogchap.com/v8go"
)

type Task struct {
	Callback *v8.Function
	Context  *v8.Context
}

// func setTimeout(info *v8.FunctionCallbackInfo) *v8.Value {
// 	var _ = make([]string, len(info.Args()))

// 	callback := info.Args()[0]

// 	if !callback.IsFunction() {
// 		fmt.Println("the first argument must be function")
// 		os.Exit(1)
// 	}
// 	v8Function, err := callback.AsFunction()
// 	if err != nil {
// 		fmt.Println("the first argument must be function")
// 		os.Exit(1)
// 	}
// 	delay := info.Args()[1].Int32()
// 	timer := time.After(time.Duration(delay) * time.Millisecond)
// 	fmt.Println("about to start timer")

// 	<-timer
// 	v8Function.Call(v8.Undefined(info.Context().Isolate()))

// 	return v8.Undefined(info.Context().Isolate())
// }

func setTimeout(info *v8.FunctionCallbackInfo) *v8.Value {

	if len(info.Args()) < 2 {
		fmt.Println("setTimeout requires at least 2 arguments")
		return v8.Undefined(info.Context().Isolate())
	}

	callback := info.Args()[0]
	if !callback.IsFunction() {
		fmt.Println("the first argument must be a function")
		return v8.Undefined(info.Context().Isolate())
	}

	function, err := callback.AsFunction()
	if err != nil {
		fmt.Println("could not convert callback to function")
		return v8.Undefined(info.Context().Isolate())
	}

	delay := info.Args()[1].Int32()
	isolate := info.Context().Isolate()
	done := make(chan bool, 1)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(delay) * time.Millisecond)
		done <- true
	}()

	wg.Wait()

	for {
		select {
		case <-done:
			function.Call(v8.Undefined(isolate))
			return v8.Undefined(isolate)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

}

var TimerFunctions = []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
}{
	{"setTimeout", func(info *v8.FunctionCallbackInfo) *v8.Value {
		go setTimeout(info)
		return v8.Undefined(info.Context().Isolate())
	}},
}

func GetTimerObjects() []struct {
	Name string
	Fn   func(*v8.FunctionCallbackInfo) *v8.Value
} {
	return TimerFunctions
}
