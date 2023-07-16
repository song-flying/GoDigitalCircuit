package digital

import (
	"context"
	"fmt"
	b "goos/bit"
)

type InWire interface {
	In() chan b.Bit
	Close()
}

type Wire struct {
	Name string
	in   chan b.Bit
	Out  chan b.Bit
}

func NewWire(ctx context.Context, name string) Wire {
	w := Wire{
		Name: name,
		in:   make(chan b.Bit, 1),
		Out:  make(chan b.Bit),
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(w.Out)
				return
			case b, ok := <-w.in:
				if !ok {
					close(w.Out)
					return
				}
				//fmt.Printf("wire %s got %s\n", w.Name, b)
				w.Out <- b
			}
		}
	}()

	return w
}

func (w *Wire) In() chan b.Bit {
	return w.in
}

func (w *Wire) Close() {
	close(w.in)
}

type MultiWire struct {
	Name  string
	in    chan b.Bit
	Wires []Wire
}

func NewMultiWire(ctx context.Context, name string, numCopies int) MultiWire {
	newName := fmt.Sprintf("%s_%d", name, numCopies)
	var wires []Wire
	for i := 0; i < numCopies; i++ {
		wires = append(wires, NewWire(ctx, fmt.Sprintf("%s_%d", newName, i)))
	}

	mw := MultiWire{
		Name:  newName,
		in:    make(chan b.Bit, 1),
		Wires: wires,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				for _, wire := range mw.Wires {
					wire.Close()
				}
				return
			case b, ok := <-mw.in:
				if !ok {
					for _, wire := range mw.Wires {
						wire.Close()
					}
					return
				}
				//fmt.Printf("multi wire %s got %s\n", mw.Name, b)
				mw.sendToAll(b)
			}
		}
	}()

	return mw
}

func (mw *MultiWire) sendToAll(b b.Bit) {
	for _, w := range mw.Wires {
		//fmt.Printf("sub wire %s got %s\n", w.Name, b)
		w.In() <- b
	}
}

func (mw *MultiWire) In() chan b.Bit {
	return mw.in
}

func (mw *MultiWire) Close() {
	close(mw.in)
}

func (w *Wire) FanOut(ctx context.Context, numCopies int) MultiWire {
	mw := NewMultiWire(ctx, w.Name, numCopies)

	go func() {
		for {
			select {
			case <-ctx.Done():
				mw.Close()
				return
			case b, ok := <-w.Out:
				if !ok {
					mw.Close()
					return
				}
				//fmt.Printf("upper wire %s got %s\n", w.Name, b)
				mw.In() <- b
			}
		}
	}()

	return mw
}
