package digital

import (
	"context"
)

type Mux struct {
	InS *MultiWire
	InA *Wire
	InB *Wire
	Out *Wire
}

func NewMux(ctx context.Context, S *MultiWire, A, B, O *Wire) *Mux {
	ns := NewWire(ctx, "~S")
	NotGate(ctx, &S.Wires[0], &ns)

	ns_and_b := NewWire(ctx, "~S&B")
	AndGate(ctx, &ns, B, &ns_and_b)

	s_and_a := NewWire(ctx, "S&A")
	AndGate(ctx, &S.Wires[1], A, &s_and_a)

	out := NewWire(ctx, "~S&B|S&A")
	OrGate(ctx, &ns_and_b, &s_and_a, &out)

	mux := Mux{
		InS: S,
		InA: A,
		InB: B,
		Out: &out,
	}

	return &mux
}
