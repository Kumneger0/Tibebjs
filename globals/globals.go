package globals

import (
	"github.com/kumneger0/tibebjs/pkg/fs"
	"github.com/kumneger0/tibebjs/pkg/net"
	v8 "rogchap.com/v8go"
)


func SetGlobalsUnderTibebNameSpace(iso *v8.Isolate) *v8.ObjectTemplate {
	globalObject := v8.NewObjectTemplate(iso)
	methods := append(fs.GetFsObjects(), net.GetNetObjects()...)

	for _, obj := range methods {
		globalObject.Set(obj.Name, v8.NewFunctionTemplate(iso, obj.Fn))
	}	
  return globalObject
}