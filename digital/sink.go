package digital

import "context"

type sink struct {
	Wire *Wire
}

func Sink(ctx context.Context, w *Wire) sink {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-w.Out:
				if !ok {
					return
				}
			}
		}
	}()

	return sink{
		Wire: w,
	}
}
