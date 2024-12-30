package fetch

import (
	"math/rand"

	"github.com/kumneger0/tibebjs/pkg/eventloop"
	v8 "rogchap.com/v8go"
)

func Fetch(info *v8.FunctionCallbackInfo) *v8.Value {

	netTask := eventloop.NetworkTask{
		Callback: &v8.Function{},
		Context:  info.Context(),
		Id:       rand.Int(),
	}

	eventloop.NetworkTaskQueue = append(eventloop.NetworkTaskQueue, netTask)
	return eventloop.Fetch(info, &netTask).Value
}
