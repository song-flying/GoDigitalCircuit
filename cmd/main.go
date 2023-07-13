package main

import (
	"context"
	"fmt"
	. "goos/bit"
	d "goos/digital"
	"goos/digital/gate"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()

	w1 := newWire(ctx, "w1")
	w2 := newWire(ctx, "w2")
	w3 := newWire(ctx, "w3")
	w4 := newWire(ctx, "w4")

	signals := [][]Bit{{O, O}, {O, I}, {I, O}, {I, I}}

	go func() {
		for _, pair := range signals {
			b1, b2 := pair[0], pair[1]
			w1.In <- b1
			w2.In <- b2
		}
		close(w1.In)
		close(w2.In)
	}()

	and := gate.AndGate(ctx, w1, w2, w3)
	not := gate.NotGate(ctx, and.Out, w4)

	var output []string
	for b := range not.Out.Out {
		output = append(output, b.String())
	}

	fmt.Printf("output = %s\n", strings.Join(output, ""))
}

func newWire(ctx context.Context, name string) d.Wire {
	return d.NewWire(ctx, name, 1*time.Second)
}
