package internal

import "reflect"

func IsMessageHandlerType(rt reflect.Type) bool {
	return rt.AssignableTo(typeOfMessageHandler)
}

func IsMessageHandler(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageHandler)
	}
	return false
}

func AsMessageHandler(rv reflect.Value) MessageHandler {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageHandler).Interface().(MessageHandler); ok {
			return v
		}
	}
	return nil
}

func isMessageObserverAffair(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageObserverAffair)
	}
	return false
}

func asMessageObserverAffair(rv reflect.Value) MessageObserverAffair {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageObserverAffair).Interface().(MessageObserverAffair); ok {
			return v
		}
	}
	return nil
}
