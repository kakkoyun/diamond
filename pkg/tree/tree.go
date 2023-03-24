package tree

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/kakkoyun/diamond/pkg"
)

type tree struct {
	ctx    context.Context
	logger log.Logger
}

func New(ctx context.Context, logger log.Logger) *tree {
	return &tree{
		ctx:    ctx,
		logger: logger,
	}
}

//go:noinline
func (t *tree) Run() error {
	ticker := time.NewTicker(pkg.RepeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.ctx.Done():
			return nil
		case <-ticker.C:
			if err := t.run(); err != nil {
				return err
			}
		}
	}
}

//go:noinline
func (t *tree) run() error {
	t.logger.Log("msg", "running tree")
	defer t.logger.Log("msg", "finished running tree")

	a(t.ctx)
	return nil
}

//go:noinline
func a(ctx context.Context) {
	b(ctx)
	b(ctx)
}

//go:noinline
func b(ctx context.Context) {
	c(ctx)
	c(ctx)
	c(ctx)
}

//go:noinline
func c(ctx context.Context) {
	leaf(ctx)
}

//go:noinline
func leaf(ctx context.Context) {
	for i := 0; i < pkg.IterCount; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
