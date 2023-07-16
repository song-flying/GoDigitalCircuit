package digital

import (
	"context"
)

type Mux struct {
	InS *DupWire
	InA *Wire
	InB *Wire
	Out *Wire
}

func NewMux(ctx context.Context, s *DupWire, a, b, o *Wire) *Mux {
	ns := NewWire(ctx, "~s")
	notS := NotGate(ctx, &s.Wire1, &ns)

	ns_and_b := NewWire(ctx, "~s&b")
	andNSB := AndGate(ctx, notS.Out, b, &ns_and_b)

	s_and_a := NewWire(ctx, "s&a")
	andSA := AndGate(ctx, &s.Wire2, a, &s_and_a)

	out := NewWire(ctx, "~s&b|s&a")
	Or := OrGate(ctx, andNSB.Out, andSA.Out, &out)

	mux := Mux{
		InS: s,
		InA: a,
		InB: b,
		Out: Or.Out,
	}

	return &mux
}
