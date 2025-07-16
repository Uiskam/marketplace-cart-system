package cqrs

import "reflect"

type Handler interface {
	Handle(command interface{}) (any, error)
}
type Processor struct {
	handlers map[reflect.Type]Handler
}

func NewProcessor() *Processor {
	return &Processor{
		handlers: make(map[reflect.Type]Handler),
	}
}

func (p *Processor) AddHandler(commandType reflect.Type, handler Handler) {
	if p.handlers == nil {
		p.handlers = make(map[reflect.Type]Handler)
	}
	p.handlers[commandType] = handler
}