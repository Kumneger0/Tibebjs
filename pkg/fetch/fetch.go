package fetch

import (
	"github.com/kumneger0/tibebjs/pkg/eventloop"
	v8 "rogchap.com/v8go"
)

func Fetch(info *v8.FunctionCallbackInfo) *v8.Value {
	return eventloop.Fetch(info).Value
}

   