package main

import (
	"context"
	"fmt"
	. "goos/bit"
	d "goos/digital"
	"strings"
	"time"
)

var D = 1 * time.Second

func main() {
	ctx := context.Background()

	signals := [][]Bit{
		{O, O, O},
		{O, O, I},
		{O, I, O},
		{O, I, I},
		{I, O, O},
		{I, O, I},
		{I, I, O},
		{I, I, I}}

	S := d.NewDupWire(ctx, "s", D)
	A := d.NewWire(ctx, "a", D)
	B := d.NewWire(ctx, "b", D)
	O := d.NewWire(ctx, "out", D)

	go func() {
		for _, t := range signals {
			s, a, b := t[0], t[1], t[2]
			fmt.Println("--------------------")
			S.In <- s
			A.In <- a
			B.In <- b
		}
		close(S.In)
		close(A.In)
		close(B.In)
	}()

	mux := d.NewMux(ctx, &S, &A, &B, &O)

	var output []string
	for b := range mux.Out.Out {
		output = append(output, b.String())
	}

	fmt.Printf("output = %s\n", strings.Join(output, ""))
}
