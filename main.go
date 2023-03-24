package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-kit/log"
	"github.com/kakkoyun/diamond/pkg/diamond"
	"github.com/kakkoyun/diamond/pkg/tree"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "name", "diamond")
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.Caller(3))

	signals := []os.Signal{os.Interrupt, os.Kill}
	ctx, cancel := context.WithCancelCause(context.Background())
	go func() {
		defer func() {
			if r := recover(); r != nil {
				cancel(fmt.Errorf("panic: %v", r))
				return
			}
			if ctx.Err() == nil {
				cancel(errors.New("signal handler exited, unexpectedly"))
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, signals...)
		select {
		case sig := <-c:
			cancel(fmt.Errorf("received signal: %v", sig))
			return
		case <-ctx.Done():
			return
		}
	}()

	go diamond.New(ctx, logger).Run()
	go tree.New(ctx, logger).Run()
	go tree.NewAnother(ctx, logger).Run()

	<-ctx.Done()
	if err := ctx.Err(); err != nil {
		logger.Log("msg", "exited with error", "err", err)
		os.Exit(1)
	}
	logger.Log("msg", "exited")
}
