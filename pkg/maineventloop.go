package mainEventloop

import (
	"sync"

	v8 "rogchap.com/v8go"
)

type Task struct {
	Callback *v8.Function
	Context  *v8.Context
}

var Tasks = make([]Task, 0)
var Wg sync.WaitGroup
var Channel = make(chan Task, 100)
