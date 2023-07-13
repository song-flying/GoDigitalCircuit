package gate

import (
	"context"
	. "goos/digital"
)

type Not struct {
	in  Wire
	Out Wire
}

func NotGate(ctx context.Context, in Wire, out Wire) *Not {
	gate := Not{
		in:  in,
		Out: out,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case b, ok := <-gate.in.Out:
				if !ok {
					gate.Out.Close()
					return
				}
				gate.Out.In <- b.Not()
			}
		}
	}()

	return &gate
}

type And struct {
	inA Wire
	inB Wire
	Out Wire
}

func AndGate(ctx context.Context, inA, inB, out Wire) *And {
	gate := And{
		inA: inA,
		inB: inB,
		Out: out,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				gate.Out.Close()
				return
			case a, ok := <-gate.inA.Out:
				if !ok {
					gate.Out.Close()
					return
				}
				select {
				case <-ctx.Done():
					return
				case b := <-gate.inB.Out:
					gate.Out.In <- a.And(b)
				}
			case b, ok := <-gate.inB.Out:
				if !ok {
					gate.Out.Close()
					return
				}
				select {
				case <-ctx.Done():
					return
				case a := <-gate.inA.Out:
					gate.Out.In <- a.And(b)
				}
			}
		}
	}()

	return &gate
}
