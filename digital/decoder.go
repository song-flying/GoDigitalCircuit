package digital

import "context"

type Decoder struct {
	A, B               *Wire
	A_B_, A_B, AB_, AB *Wire
}

func NewDecoder(ctx context.Context, A, B *Wire, A_B_, A_B, AB_, AB *Wire) Decoder {
	A3 := A.FanOut(ctx, 3)
	B3 := B.FanOut(ctx, 3)

	A_ := NewWire(ctx, "~A")
	NotGate(ctx, &A3.Wires[0], &A_)
	A_2 := A_.FanOut(ctx, 2)

	B_ := NewWire(ctx, "~B")
	NotGate(ctx, &B3.Wires[0], &B_)
	B_2 := B_.FanOut(ctx, 2)

	AndGate(ctx, &A_2.Wires[0], &B_2.Wires[0], A_B_)
	AndGate(ctx, &A_2.Wires[1], &B3.Wires[1], A_B)
	AndGate(ctx, &A3.Wires[1], &B_2.Wires[1], AB_)
	AndGate(ctx, &A3.Wires[2], &B3.Wires[2], AB)

	d := Decoder{
		A:    A,
		B:    B,
		A_B_: A_B_,
		A_B:  A_B,
		AB_:  AB_,
		AB:   AB,
	}

	return d
}
