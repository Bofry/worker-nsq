package internal

import "reflect"

type Router map[string]RouteComponent

func (r Router) Add(topic string, handler MessageHandler, handlerComponentID string) {
	r[topic] = RouteComponent{
		MessageHandler:     handler,
		HandlerComponentID: handlerComponentID,
		HandlerType:        reflect.TypeOf(handler),
	}
}

func (r Router) Remove(topic string) {
	delete(r, topic)
}

func (r Router) Get(topic string) MessageHandler {
	if r == nil {
		return nil
	}

	if v, ok := r[topic]; ok {
		return v.MessageHandler
	}
	return nil
}

func (r Router) Has(topic string) bool {
	if r == nil {
		return false
	}

	if _, ok := r[topic]; ok {
		return true
	}
	return false
}

func (r Router) FindHandlerComponentID(topic string) string {
	if r == nil {
		return ""
	}

	if v, ok := r[topic]; ok {
		return v.HandlerComponentID
	}
	return ""
}

func (r Router) FindHandlerType(topic string) reflect.Type {
	if r == nil {
		return nil
	}

	if v, ok := r[topic]; ok {
		return v.HandlerType
	}
	return nil
}
