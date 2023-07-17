package digital

import (
	"context"
	"goos/bit"
)

type TrafficLight struct {
	Clock            *Wire
	EWCar, NSCar     *Wire
	EWLight, NSLight *Wire
	state            FlipFlop
}

func NewTrafficLight(ctx context.Context, Clock, EWCar, NSCar, EWLight, NSLight *Wire) TrafficLight {
	EWNext := NewWire(ctx, "~EW*EWCar+EW*~NSCar")
	C := Clock
	D := &EWNext
	Q := NewMultiWire(ctx, "EW", 2)
	Q_ := NewMultiWire(ctx, "~EW", 2)
	ff := NewFlipFlop(ctx, C, D, &Q, &Q_)

	// bootstrap
	D.In() <- bit.O

	EW := ff.Q.FanOut(ctx, 2)
	EW_ := ff.Q_.FanOut(ctx, 2)

	NSCar_ := NewWire(ctx, "~NSCar")
	NotGate(ctx, NSCar, &NSCar_)

	AND1 := NewWire(ctx, "~EW*EWCar")
	AndGate(ctx, &EW_.Wires[0], EWCar, &AND1)

	AND2 := NewWire(ctx, "EW*~NSCar")
	AndGate(ctx, &EW.Wires[0], &NSCar_, &AND2)

	OrGate(ctx, &AND1, &AND2, &EWNext)

	IDGate(ctx, &EW.Wires[1], EWLight)
	IDGate(ctx, &EW_.Wires[1], NSLight)

	return TrafficLight{
		Clock:   Clock,
		EWCar:   EWCar,
		NSCar:   NSCar,
		EWLight: EWLight,
		NSLight: NSLight,
		state:   ff,
	}
}
