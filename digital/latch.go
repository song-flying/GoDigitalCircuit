package digital

import (
	"context"
	"goos/bit"
)

type SRLatch struct {
	R, S  *Wire
	Q, Q_ *Wire
}

func NewSRLatch(ctx context.Context, R, S *Wire, Q, Q_ *MultiWire) SRLatch {
	NorGate(ctx, R, &Q_.Wires[1], Q)
	NorGate(ctx, S, &Q.Wires[1], Q_)

	// bootstrap
	Q.Wires[1].In() <- bit.O
	Q_.Wires[1].In() <- bit.I

	return SRLatch{
		R:  R,
		S:  S,
		Q:  &Q.Wires[0],
		Q_: &Q_.Wires[0],
	}
}

type DLatch struct {
	sr    SRLatch
	C, D  *Wire
	Q, Q_ *Wire
}

func NewDLatch(ctx context.Context, C, D *Wire, Q, Q_ *MultiWire) DLatch {
	C2 := C.FanOut(ctx, 2)
	D2 := D.FanOut(ctx, 2)
	D_ := NewWire(ctx, "~D")
	NotGate(ctx, &D2.Wires[0], &D_)

	R := NewWire(ctx, "C&~D")
	S := NewWire(ctx, "C&D")
	AndGate(ctx, &C2.Wires[0], &D_, &R)
	AndGate(ctx, &C2.Wires[1], &D2.Wires[1], &S)

	sr := NewSRLatch(ctx, &R, &S, Q, Q_)

	return DLatch{
		sr: sr,
		C:  C,
		D:  D,
		Q:  &Q.Wires[0],
		Q_: &Q_.Wires[0],
	}
}
