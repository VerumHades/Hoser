package outbox

import (
	"common/pkg/shared"
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"
)

type Dispatcher struct {
	queryRepository   OutboxQueryRepository
	commandRepository OutboxCommandRepository

	pollInterval time.Duration

	cursor OutboxEventCursor

	listeners map[string][]EventListener
	mutex     sync.RWMutex
}

func NewDispatcher(
	queryRepository OutboxQueryRepository,
	commandRepository OutboxCommandRepository,
	initialCursor OutboxEventCursor,
	pollInterval time.Duration,
) *Dispatcher {
	return &Dispatcher{
		queryRepository:   queryRepository,
		commandRepository: commandRepository,
		pollInterval:      pollInterval,
		cursor:            initialCursor,
		listeners:         make(map[string][]EventListener),
	}
}

func (dispatcher *Dispatcher) Run(
	ctx context.Context,
) error {
	ticker := time.NewTicker(dispatcher.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := dispatcher.dispatchOnce(ctx); err != nil {
				return err
			}
		}
	}
}

func (dispatcher *Dispatcher) dispatchOnce(
	ctx context.Context,
) error {
	events, nextCursor, err := dispatcher.queryRepository.FetchNextBatchAfter(
		ctx,
		dispatcher.cursor.LastOccurredAt,
		shared.BatchRequest[OutboxEventCursor]{
			Cursor:       dispatcher.cursor,
			MaxBatchSize: 50,
		},
	)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := dispatcher.dispatchEvent(ctx, event); err != nil {
			return err
		}
	}

	dispatcher.cursor = nextCursor
	return nil
}

func (dispatcher *Dispatcher) dispatchEvent(
	ctx context.Context,
	event *OutboxEvent,
) error {
	dispatcher.mutex.RLock()
	listeners := dispatcher.listeners[event.eventType]
	dispatcher.mutex.RUnlock()

	for _, listener := range listeners {
		if err := listener.Handle(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func (dispatcher *Dispatcher) RegisterListener(
	eventType string,
	listener EventListener,
) {
	dispatcher.mutex.Lock()
	defer dispatcher.mutex.Unlock()

	dispatcher.listeners[eventType] = append(
		dispatcher.listeners[eventType],
		listener,
	)
}

type EventListener interface {
	Handle(
		ctx context.Context,
		event *OutboxEvent,
	) error
}

type TypedEventListener[T any] struct {
	handler func(ctx context.Context, payload T) error
}

func NewTypedEventListener[T any](
	handler func(ctx context.Context, payload T) error,
) *TypedEventListener[T] {
	return &TypedEventListener[T]{
		handler: handler,
	}
}

func (listener *TypedEventListener[T]) Handle(
	ctx context.Context,
	event *OutboxEvent,
) error {
	var payload T

	if err := json.Unmarshal(event.payload, &payload); err != nil {
		return err
	}

	return listener.handler(ctx, payload)
}

func RegisterTypedListener[T any](
	dispatcher *Dispatcher,
	handler func(ctx context.Context, payload T) error,
) {
	eventType := EventTypeFromPayload[T]()

	dispatcher.RegisterListener(
		eventType,
		NewTypedEventListener(handler),
	)
}

func EventTypeFromPayload[T any]() string {
	var payload T
	payloadType := reflect.TypeOf(payload)

	if payloadType.Kind() == reflect.Pointer {
		payloadType = payloadType.Elem()
	}

	return payloadType.Name()
}
