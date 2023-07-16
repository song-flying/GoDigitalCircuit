package main

import (
	"context"
	"fmt"
	"goos/bit"
	d "goos/digital"
	"strings"
	"time"
)

func main() {
	//fmt.Println("--------------- Test Mux ---------------")
	//testMux()
	fmt.Println("--------------- Test Decoder ---------------")
	testDecoder()
	//fmt.Println("--------------- Test SRLatch ---------------")
	//testSRLatch()
}

func testMux() {
	ctx := context.Background()

	signals := [][]bit.Bit{
		{0, 0, 0},
		{0, 0, 1},
		{0, 1, 0},
		{0, 1, 1},
		{1, 0, 0},
		{1, 0, 1},
		{1, 1, 0},
		{1, 1, 1}}

	S := d.NewDupWire(ctx, "s")
	A := d.NewWire(ctx, "a")
	B := d.NewWire(ctx, "b")
	O := d.NewWire(ctx, "out")

	signalCh := make(chan []bit.Bit)
	source := d.NewSource(ctx, signalCh, []chan bit.Bit{S.In, A.In, B.In})

	go func() {
		for _, v := range signals {
			source.C <- v
			time.Sleep(2 * time.Second)
		}
		close(source.C)
	}()

	mux := d.NewMux(ctx, &S, &A, &B, &O)

	var output []string
	for b := range mux.Out.Out {
		output = append(output, fmt.Sprintf("%3s", b.String()))
		fmt.Printf("\routput = %s", strings.Join(output, "|"))
	}

	fmt.Println()
}

func testDecoder() {
	ctx := context.Background()

	A := d.NewWire(ctx, "A")
	B := d.NewWire(ctx, "B")
	A_B_ := d.NewWire(ctx, "~A~B")
	A_B := d.NewWire(ctx, "~AB")
	AB_ := d.NewWire(ctx, "A~B")
	AB := d.NewWire(ctx, "AB")

	values := [][]bit.Bit{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	}

	signalCh := make(chan []bit.Bit)
	source := d.NewSource(ctx, signalCh, []chan bit.Bit{A.In, B.In})

	go func() {
		for _, v := range values {
			source.C <- v
			time.Sleep(3 * time.Second)
		}
		close(source.C)
	}()

	d := d.NewDecoder(ctx, &A, &B, &A_B_, &A_B, &AB_, &AB)

	var output []string
	for {
		out1, ok1 := <-d.A_B_.Out
		out2, ok2 := <-d.A_B.Out
		out3, ok3 := <-d.AB_.Out
		out4, ok4 := <-d.AB.Out
		if !ok1 || !ok2 || !ok3 || !ok4 {
			break
		}
		output = append(output, fmt.Sprintf("%s%s%s%s", out1, out2, out3, out4))
		fmt.Printf("\routput = %s", strings.Join(output, "|"))
	}
	fmt.Println()
}

func testSRLatch() {
	ctx := context.Background()

	R := d.NewWire(ctx, "r")
	S := d.NewWire(ctx, "s")
	Q := d.NewDupWire(ctx, "q")
	Q_ := d.NewDupWire(ctx, "~q")

	values := [][]bit.Bit{
		{0, 0},
		{0, 1},
		{0, 0},
		{1, 0},
	}
	signalCh := make(chan []bit.Bit)
	source := d.NewSource(ctx, signalCh, []chan bit.Bit{R.In, S.In})
	go func() {
		for _, v := range values {
			source.C <- v
			time.Sleep(5 * time.Second)
		}
		close(source.C)
	}()

	latch := d.NewSRLatch(ctx, &R, &S, &Q, &Q_)

	var output []string
	for {
		q, ok1 := <-latch.Q.Out
		q_, ok2 := <-latch.Q_.Out
		if !ok1 || !ok2 {
			break
		}
		output = append(output, fmt.Sprintf("%s%s", q, q_))
		fmt.Printf("\routput = %s", strings.Join(output, "|"))
	}
	fmt.Println()
}
