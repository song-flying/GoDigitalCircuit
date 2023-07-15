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
	fmt.Println("----------Test Mux----------")
	testMux()
	fmt.Println("----------Test SRLatch----------")
	testSRLatch()
}

func testSRLatch() {
	ctx := context.Background()

	R := d.NewWire(ctx, "r", D)
	S := d.NewWire(ctx, "s", D)
	Q := d.NewDupWire(ctx, "q", D)
	Q_ := d.NewDupWire(ctx, "-q", D)

	signals := [][]Bit{
		{0, 0},
		{0, 1},
		{0, 1},
		{0, 0},
		{1, 0},
		{1, 0},
		{0, 0},
	}
	go func() {
		for _, t := range signals {
			r, s := t[0], t[1]
			R.In <- r
			S.In <- s
		}
		close(R.In)
		close(S.In)
	}()

	latch := d.NewSRLatch(ctx, &R, &S, &Q, &Q_)

	var output []string
	for {
		q, ok1 := <-latch.Q.Out
		q_, ok2 := <-latch.Q_.Out
		if !ok1 || !ok2 {
			break
		}
		output = append(output, fmt.Sprintf("(%s,%s)", q, q_))
	}

	fmt.Printf("output = %s\n", strings.Join(output, ""))
}

func testMux() {
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
