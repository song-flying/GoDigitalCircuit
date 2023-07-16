package digital

import "context"

type FlipFlop struct {
	master DLatch
	slave  DLatch
	C, D   *Wire
	Q, Q_  *Wire
}

func NewFlipFlop(ctx context.Context, C, D *Wire, Q, Q_ *MultiWire) FlipFlop {
	C2 := C.FanOut(ctx, 2)

	MD := D
	MC := C2.Wires[0]
	MQ := NewMultiWire(ctx, "Q", 2)
	MQ_ := NewMultiWire(ctx, "Q_", 2)
	md := NewDLatch(ctx, &MC, MD, &MQ, &MQ_)
	Sink(ctx, md.Q_)

	C_ := NewWire(ctx, "~C")
	NotGate(ctx, &C2.Wires[1], &C_)

	SD := md.Q
	SC := C_
	SQ := Q
	SQ_ := Q_

	sd := NewDLatch(ctx, &SC, SD, SQ, SQ_)

	return FlipFlop{
		master: md,
		slave:  sd,
		C:      C,
		D:      D,
		Q:      sd.Q,
		Q_:     sd.Q_,
	}
}
