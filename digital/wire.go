package digital

import (
	"context"
	"fmt"
	b "goos/bit"
)

type Wire struct {
	Name string
	In   chan b.Bit
	Out  chan b.Bit
}

func NewWire(ctx context.Context, name string) Wire {
	w := Wire{
		Name: name,
		In:   make(chan b.Bit, 1),
		Out:  make(chan b.Bit),
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(w.Out)
				return
			case b, ok := <-w.In:
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

func (w *Wire) Close() {
	close(w.In)
}

type DupWire struct {
	name  string
	In    chan b.Bit
	Wire1 Wire
	Wire2 Wire
}

func NewDupWire(ctx context.Context, name string) DupWire {
	w := DupWire{
		name:  name,
		In:    make(chan b.Bit, 1),
		Wire1: NewWire(ctx, name),
		Wire2: NewWire(ctx, name),
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				w.Wire1.Close()
				w.Wire2.Close()
				return
			case b, ok := <-w.In:
				if !ok {
					w.Wire1.Close()
					w.Wire2.Close()
					return
				}
				w.Wire1.In <- b
				w.Wire2.In <- b
			}
		}
	}()

	return w
}

func (dw *DupWire) Close() {
	close(dw.In)
}

type MultiWire struct {
	Name  string
	In    chan b.Bit
	Wires []Wire
}

func NewMultiWire(ctx context.Context, name string, size int) MultiWire {
	newName := fmt.Sprintf("%s_%d", name, size)
	var wires []Wire
	for i := 0; i < size; i++ {
		wires = append(wires, NewWire(ctx, fmt.Sprintf("%s_%d", newName, i)))
	}

	mw := MultiWire{
		Name:  newName,
		In:    make(chan b.Bit, 1),
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
			case b, ok := <-mw.In:
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
		w.In <- b
	}
}

func (mw *MultiWire) Close() {
	close(mw.In)
}

func (w *Wire) FanOut(ctx context.Context, size int) MultiWire {
	mw := NewMultiWire(ctx, w.Name, size)

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
				mw.In <- b
			}
		}
	}()

	return mw
}
