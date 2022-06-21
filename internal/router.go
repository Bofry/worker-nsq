package internal

type Router map[string]MessageHandler

func NewRouter() *Router {
	var router Router = make(map[string]MessageHandler)
	return &router
}

func (r *Router) Add(topic string, handler MessageHandler) {
	route := *r
	route[topic] = handler
}

func (r *Router) Remove(topic string) {
	route := *r
	delete(route, topic)
}

func (r *Router) Get(topic string) MessageHandler {
	route := *r
	if v, ok := route[topic]; ok {
		return v
	}
	return nil
}
