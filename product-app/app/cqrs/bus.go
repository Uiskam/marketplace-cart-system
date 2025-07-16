package cqrs

import (
	"fmt"
	"reflect"
)

type Bus struct {
	Processor *Processor
}

func NewBus(processor *Processor) *Bus {
	return &Bus{
		Processor: processor,
	}
}

func (b *Bus) Send(command interface{}) (any, error) {
	commandType := reflect.TypeOf(command)
	handler, exists := b.Processor.handlers[commandType]
	if !exists {
		return nil, fmt.Errorf("no handler found for command type %s", commandType)
	}
	result, err := handler.Handle(command)
	if err != nil {
		return nil, fmt.Errorf("error handling command %s: %w", commandType, err)
	}
	return result, nil
}