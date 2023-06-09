package internal

import (
	"reflect"
	"sync"

	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

type MessageTracerService struct {
	TracerProvider    *trace.SeverityTracerProvider
	TextMapPropagator propagation.TextMapPropagator

	Enabled bool

	UnhandledMessageHandlerComponentID string

	unhandledMessageTracer *trace.SeverityTracer

	tracers            map[string]*trace.SeverityTracer
	tracersInitializer sync.Once
}

func NewMessageTracerService() *MessageTracerService {
	return &MessageTracerService{}
}

func (s *MessageTracerService) Tracer(id string) *trace.SeverityTracer {
	if s.tracers != nil {
		if tr, ok := s.tracers[id]; ok {
			return tr
		}
	}
	return s.unhandledMessageTracer
}

func (s *MessageTracerService) init(messageManager interface{}) {
	if s.TextMapPropagator == nil {
		s.TextMapPropagator = defaultTextMapPropagator
	}
	if s.TracerProvider == nil {
		s.TracerProvider = defaultTracerProvider
	}
	if s.Enabled {
		s.makeTracerMap()
		s.buildTracer(messageManager)
	}
	s.makeUnhandledMessageTracer()
}

func (s *MessageTracerService) makeTracerMap() {
	s.tracersInitializer.Do(func() {
		s.tracers = make(map[string]*trace.SeverityTracer)
	})
}

func (s *MessageTracerService) buildTracer(requestManager interface{}) {
	var (
		rvManager reflect.Value = reflect.ValueOf(requestManager)
	)
	if rvManager.Kind() != reflect.Pointer || rvManager.IsNil() {
		return
	}

	rvManager = reflect.Indirect(rvManager)
	numOfHandles := rvManager.NumField()
	for i := 0; i < numOfHandles; i++ {
		rvRequest := rvManager.Field(i)
		if rvRequest.Kind() != reflect.Pointer || rvRequest.IsNil() {
			continue
		}

		rvRequest = reflect.Indirect(rvRequest)
		if rvRequest.Kind() == reflect.Struct {
			rvRequest = reflect.Indirect(rvRequest)

			componentName := rvRequest.Type().Name()
			tracer := s.TracerProvider.Tracer(componentName)

			info := rvManager.Type().Field(i)
			if _, ok := s.tracers[info.Name]; !ok {
				s.registerTracer(info.Name, tracer)
			}
		}
	}
}

func (s *MessageTracerService) registerTracer(id string, tracer *trace.SeverityTracer) {
	container := s.tracers

	if tracer != nil {
		if _, ok := container[id]; ok {
			NsqWorkerLogger.Fatalf("specified id '%s' already exists", id)
		}
		container[id] = tracer
	}
}

func (s *MessageTracerService) makeUnhandledMessageTracer() {
	var (
		tp *trace.SeverityTracerProvider = defaultTracerProvider
	)

	if s.Enabled {
		tp = s.TracerProvider
	}

	if len(s.UnhandledMessageHandlerComponentID) > 0 {
		v, ok := s.tracers[s.UnhandledMessageHandlerComponentID]
		if ok {
			s.unhandledMessageTracer = v
			return
		}
	}
	s.unhandledMessageTracer = tp.Tracer("")
}
