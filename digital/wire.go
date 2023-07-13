package digital

import (
	"context"
	"fmt"
	b "goos/bit"
	"time"
)

type Wire struct {
	name string
	In   chan b.Bit
	Out  chan b.Bit
}

func NewWire(ctx context.Context, name string, duration time.Duration) Wire {
	w := Wire{
		name: name,
		In:   make(chan b.Bit),
		Out:  make(chan b.Bit),
	}

	go func() {
		for {
			select {
			case b, ok := <-w.In:
				if !ok {
					close(w.Out)
					return
				}
				fmt.Printf("%s <-- %s\n", w.name, b)
				if duration != 0 {
					time.Sleep(duration)
				}
				w.Out <- b
			case <-ctx.Done():
				close(w.Out)
				return
			}
		}
	}()

	return w
}

func (w *Wire) Close() {
	close(w.In)
}
