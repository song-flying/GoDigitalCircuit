package digital

import (
	"context"
	b "goos/bit"
)

type UnaryGate struct {
	in  *Wire
	Out *Wire
}

func NotGate(ctx context.Context, in *Wire, out *Wire) UnaryGate {
	gate := UnaryGate{
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

	return gate
}

type BinaryGate struct {
	inA *Wire
	inB *Wire
	Out *Wire
}

func binaryGate(ctx context.Context, inA, inB, out *Wire, f func(b.Bit, b.Bit) b.Bit) BinaryGate {
	gate := BinaryGate{
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
					gate.Out.In <- f(a, b)
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
					gate.Out.In <- f(a, b)
				}
			}
		}
	}()

	return gate
}

func AndGate(ctx context.Context, inA, inB, out *Wire) BinaryGate {
	return binaryGate(ctx, inA, inB, out, b.And)
}

func OrGate(ctx context.Context, inA, inB, out *Wire) BinaryGate {
	return binaryGate(ctx, inA, inB, out, b.Or)
}

type BinaryDupGate struct {
	inA *Wire
	inB *Wire
	Out *DupWire
}

func binaryDupGate(ctx context.Context, inA, inB *Wire, out *DupWire, f func(b.Bit, b.Bit) b.Bit) BinaryDupGate {
	gate := BinaryDupGate{
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
					gate.Out.In <- f(a, b)
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
					gate.Out.In <- f(a, b)
				}
			}
		}
	}()

	return gate
}

func NORGate(ctx context.Context, inA, inB *Wire, out *DupWire) BinaryDupGate {
	return binaryDupGate(ctx, inA, inB, out, b.Nor)
}
