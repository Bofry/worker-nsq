package internal

import "reflect"

type RouteComponent struct {
	MessageHandler     MessageHandler
	HandlerComponentID string
	HandlerType        reflect.Type
}
