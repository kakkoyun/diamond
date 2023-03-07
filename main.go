package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/go-kit/log"
	"github.com/metalmatze/signal/healthcheck"
	"github.com/metalmatze/signal/internalserver"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Flags struct {
	Address string `default:":8080" help:"Address string for internal server"`
}

func main() {
	flags := &Flags{}
	_ = kong.Parse(flags)

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "name", "diamond")
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.Caller(3))

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewBuildInfoCollector(),
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	healthchecks := healthcheck.NewMetricsHandler(healthcheck.NewHandler(), registry)
	h := internalserver.NewHandler(
		internalserver.WithHealthchecks(healthchecks),
		internalserver.WithPrometheusRegistry(registry),
		internalserver.WithPProf(),
	)

	s := http.Server{
		Addr:    flags.Address,
		Handler: h,
	}

	var g run.Group

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	{
		g.Add(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					a()
					b()
				}
			}
		}, func(err error) {
			cancel()
		})
	}

	g.Add(func() error {
		logger.Log("msg", "starting internal HTTP server", "address", s.Addr)
		return s.ListenAndServe()
	}, func(err error) {
		_ = s.Shutdown(context.Background())
	})

	g.Add(run.SignalHandler(ctx, os.Interrupt, os.Kill))
	if err := g.Run(); err != nil {
		var e run.SignalError
		if errors.As(err, &e) {
			logger.Log("msg", "program exited with signal", "err", err, "signal", e.Signal)
		} else {
			logger.Log("msg", "program exited with error", "err", err)
		}
		os.Exit(1)
	}
	logger.Log("msg", "exited")
}

//go:noinline
func a() {
	println("a")
	c()
}

//go:noinline
func b() {
	println("b")
	c()
}

//go:noinline
func c() {
	println("c")
}
