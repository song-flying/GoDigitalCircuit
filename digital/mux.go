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

func NewMux(ctx context.Context, s *MultiWire, a, b, o *Wire) *Mux {
	ns := NewWire(ctx, "~s")
	NotGate(ctx, &s.Wires[0], &ns)

	ns_and_b := NewWire(ctx, "~s&b")
	AndGate(ctx, &ns, b, &ns_and_b)

	s_and_a := NewWire(ctx, "s&a")
	AndGate(ctx, &s.Wires[1], a, &s_and_a)

	out := NewWire(ctx, "~s&b|s&a")
	OrGate(ctx, &ns_and_b, &s_and_a, &out)

	mux := Mux{
		InS: s,
		InA: a,
		InB: b,
		Out: &out,
	}

	return &mux
}
