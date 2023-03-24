package tree

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/kakkoyun/diamond/pkg"
)

type anotherTree struct {
	ctx    context.Context
	logger log.Logger
}

func NewAnother(ctx context.Context, logger log.Logger) *anotherTree {
	return &anotherTree{
		ctx:    ctx,
		logger: logger,
	}
}

//go:noinline
func (t *anotherTree) Run() error {
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
func (t *anotherTree) run() error {
	t.logger.Log("msg", "running another tree")
	defer t.logger.Log("msg", "finished running another tree")

	d(t.ctx)
	return nil
}

//go:noinline
func d(ctx context.Context) {
	e(ctx)
	f(ctx)
}

//go:noinline
func e(ctx context.Context) {
	// g(ctx)
	// g(ctx)
	// g(ctx)
	h(ctx)
}

//go:noinline
func f(ctx context.Context) {
	h(ctx)
	h(ctx)
	h(ctx)
}

//go:noinline
func g(ctx context.Context) {
	leaf(ctx)
}

//go:noinline
func h(ctx context.Context) {
	leaf(ctx)
}
