package tree

import "context"

//go:noinline
func Run(ctx context.Context) {
	a(ctx)
}

//go:noinline
func c(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

//go:noinline
func b(ctx context.Context) {
	c(ctx)
	c(ctx)
	c(ctx)
}

//go:noinline
func a(ctx context.Context) {
	b(ctx)
	b(ctx)
}
