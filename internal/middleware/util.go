package middleware

import (
	"reflect"
	"unsafe"

	"github.com/Bofry/worker-nsq/internal"
)

func isMessageHandler(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageHandler)
	}
	return false
}

func isMessageObserver(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageObserver)
	}
	return false
}

func asMessageHandler(rv reflect.Value) internal.MessageHandler {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageHandler).Interface().(internal.MessageHandler); ok {
			return v
		}
	}
	return nil
}

func asMessageObserver(rv reflect.Value) internal.MessageObserver {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageObserver).Interface().(internal.MessageObserver); ok {
			return v
		}
	}
	return nil
}

func asNsqWorker(rv reflect.Value) *internal.NsqWorker {
	return reflect.NewAt(typeOfHost, unsafe.Pointer(rv.Pointer())).
		Interface().(*internal.NsqWorker)
}
