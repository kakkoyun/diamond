package diamond

import (
	"context"
	"time"

	"github.com/go-kit/log"

	"github.com/kakkoyun/diamond/pkg"
)

type diamond struct {
	ctx    context.Context
	logger log.Logger
}

func New(ctx context.Context, logger log.Logger) *diamond {
	return &diamond{
		ctx:    ctx,
		logger: logger,
	}
}

//go:noinline
func (d *diamond) Run() error {
	ticker := time.NewTicker(pkg.RepeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return nil
		case <-ticker.C:
			if err := d.run(); err != nil {
				return err
			}
		}
	}
}

//go:noinline
func (d *diamond) run() error {
	d.logger.Log("msg", "running diamond")
	defer d.logger.Log("msg", "finished running diamond")

	a(d.ctx)
	b(d.ctx)
	return nil
}

//go:noinline
func a(ctx context.Context) {
	c(ctx)
}

//go:noinline
func b(ctx context.Context) {
	c(ctx)
}

//go:noinline
func c(ctx context.Context) {
	for i := 0; i < pkg.IterCount; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
