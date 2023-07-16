package digital

import (
	"context"
	"goos/bit"
)

type SRLatch struct {
	R, S  *Wire
	Q, Q_ *Wire
}

func NewSRLatch(ctx context.Context, r, s *Wire, q, q_ *MultiWire) SRLatch {
	NORGate(ctx, r, &q_.Wires[1], q)
	NORGate(ctx, s, &q.Wires[1], q_)

	// bootstrap
	q.Wires[1].In() <- bit.O
	q_.Wires[1].In() <- bit.I

	return SRLatch{
		R:  r,
		S:  s,
		Q:  &q.Wires[0],
		Q_: &q_.Wires[0],
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
