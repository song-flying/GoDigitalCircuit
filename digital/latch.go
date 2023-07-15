package digital

import (
	"context"
	"goos/bit"
)

type SRLatch struct {
	R  *Wire
	S  *Wire
	Q  *Wire
	Q_ *Wire
}

func NewSRLatch(ctx context.Context, r, s *Wire, q, q_ *DupWire) *SRLatch {
	NORGate(ctx, r, &q_.Wire2, q)
	NORGate(ctx, s, &q.Wire2, q_)

	// bootstrap
	q.Wire2.In <- bit.O
	q_.Wire2.In <- bit.I

	return &SRLatch{
		R:  r,
		S:  s,
		Q:  &q.Wire1,
		Q_: &q_.Wire1,
	}
}
