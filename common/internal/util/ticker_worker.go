package util

import (
	"context"
	"time"
)

type TickFunction func(context.Context) error
type ErrorHandler func(error)

type TickerWorker struct {
	interval     time.Duration
	tickFunction TickFunction
	errorHandler ErrorHandler
	cancelFunc   context.CancelFunc
}

func NewTickerWorker(interval time.Duration, tickFunction TickFunction, errorHandler ErrorHandler) *TickerWorker {
	if tickFunction == nil {
		panic("tickFunction must not be nil")
	}

	if errorHandler == nil {
		errorHandler = func(error) {}
	}

	return &TickerWorker{
		interval:     interval,
		tickFunction: tickFunction,
		errorHandler: errorHandler,
	}
}

func (tickerWorker *TickerWorker) Start(parentContext context.Context) {
	if tickerWorker.cancelFunc != nil {
		panic("TickerWorker already started")
	}

	workerContext, cancelFunc := context.WithCancel(parentContext)
	tickerWorker.cancelFunc = cancelFunc

	go func() {
		ticker := time.NewTicker(tickerWorker.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := tickerWorker.tickFunction(workerContext); err != nil {
					tickerWorker.errorHandler(err)
				}
			case <-workerContext.Done():
				return
			}
		}
	}()
}

func (tickerWorker *TickerWorker) Stop() {
	if tickerWorker.cancelFunc != nil {
		tickerWorker.cancelFunc()
		tickerWorker.cancelFunc = nil
	}
}
