package diamond

import "context"

//go:noinline
func Run(ctx context.Context) {
	a(ctx)
	b(ctx)
}

//go:noinline
func a(ctx context.Context) {
	println("a")
	c(ctx)
}

//go:noinline
func b(ctx context.Context) {
	println("b")
	c(ctx)
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
