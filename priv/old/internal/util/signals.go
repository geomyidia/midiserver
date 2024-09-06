package util

import (
	"context"
	"os"
	"os/signal"
)

func HandleSignal(handler func(int, os.Signal), signals ...os.Signal) {
	signalHandler := make(chan os.Signal, 1)
	signal.Notify(signalHandler, signals...)
	s := <-signalHandler
	handler(os.Getpid(), s)
}

func SignalWithContext(ctx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, signals...)
}
