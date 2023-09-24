package internal

import "reflect"

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
