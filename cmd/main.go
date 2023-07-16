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
	//fmt.Println("--------------- Test Decoder ---------------")
	//testDecoder()
	//fmt.Println("--------------- Test SRLatch ---------------")
	//testSRLatch()
	//fmt.Println("--------------- Test DLatch ---------------")
	//testDLatch()
	fmt.Println("--------------- Test FlipFlop ---------------")
	testFlipFlop()
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

	S := d.NewMultiWire(ctx, "S", 2)
	A := d.NewWire(ctx, "A")
	B := d.NewWire(ctx, "B")
	O := d.NewWire(ctx, "O")

	signalCh := make(chan []bit.Bit)
	source := d.NewSource(ctx, signalCh, []chan bit.Bit{S.In(), A.In(), B.In()})

	go func() {
		for _, v := range signals {
			source.C <- v
			time.Sleep(2 * time.Second)
		}
		close(source.C)
	}()

	m := d.NewMux(ctx, &S, &A, &B, &O)

	var output []string
	for b := range m.Out.Out {
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
	source := d.NewSource(ctx, signalCh, []chan bit.Bit{A.In(), B.In()})

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

	R := d.NewWire(ctx, "R")
	S := d.NewWire(ctx, "S")
	Q := d.NewMultiWire(ctx, "Q", 2)
	Q_ := d.NewMultiWire(ctx, "~Q", 2)

	values := [][]bit.Bit{
		{0, 0},
		{0, 1},
		{0, 0},
		{1, 0},
	}
	signalCh := make(chan []bit.Bit)
	source := d.NewSource(ctx, signalCh, []chan bit.Bit{R.In(), S.In()})
	go func() {
		for _, v := range values {
			source.C <- v
			time.Sleep(5 * time.Second)
		}
		close(source.C)
	}()

	sr := d.NewSRLatch(ctx, &R, &S, &Q, &Q_)

	var output []string
	for {
		q, ok1 := <-sr.Q.Out
		q_, ok2 := <-sr.Q_.Out
		if !ok1 || !ok2 {
			break
		}
		output = append(output, fmt.Sprintf("%s%s", q, q_))
		fmt.Printf("\routput = %s", strings.Join(output, "|"))
	}
	fmt.Println()
}

func testDLatch() {
	ctx := context.Background()

	C := d.NewWire(ctx, "C")
	D := d.NewWire(ctx, "D")
	Q := d.NewMultiWire(ctx, "Q", 2)
	Q_ := d.NewMultiWire(ctx, "~Q", 2)

	values := [][]bit.Bit{
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 0},
		{0, 1},
		{0, 0},
		{0, 1},
	}
	signalCh := make(chan []bit.Bit)
	source := d.NewSource(ctx, signalCh, []chan bit.Bit{C.In(), D.In()})
	go func() {
		for _, v := range values {
			source.C <- v
			time.Sleep(5 * time.Second)
		}
		close(source.C)
	}()

	d := d.NewDLatch(ctx, &C, &D, &Q, &Q_)

	var output []string
	for {
		q, ok1 := <-d.Q.Out
		q_, ok2 := <-d.Q_.Out
		if !ok1 || !ok2 {
			break
		}
		output = append(output, fmt.Sprintf("%s%s", q, q_))
		//fmt.Printf("\routput = %s", strings.Join(output, "|"))
		fmt.Printf("output = %s\n", strings.Join(output, "|"))
	}
	fmt.Println()
}

func testFlipFlop() {
	ctx := context.Background()

	C := d.NewWire(ctx, "C")
	D := d.NewWire(ctx, "D")
	Q := d.NewMultiWire(ctx, "Q", 2)
	Q_ := d.NewMultiWire(ctx, "~Q", 2)

	values := [][]bit.Bit{
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 1},
		{0, 0},
		{0, 0},
		{0, 1},
	}
	signalCh := make(chan []bit.Bit)
	source := d.NewSource(ctx, signalCh, []chan bit.Bit{C.In(), D.In()})
	go func() {
		for _, v := range values {
			source.C <- v
			time.Sleep(5 * time.Second)
		}
		close(source.C)
	}()

	d := d.NewFlipFlop(ctx, &C, &D, &Q, &Q_)

	var output []string
	for {
		q, ok1 := <-d.Q.Out
		q_, ok2 := <-d.Q_.Out
		if !ok1 || !ok2 {
			break
		}
		output = append(output, fmt.Sprintf("%s%s", q, q_))
		//fmt.Printf("\routput = %s", strings.Join(output, "|"))
		fmt.Printf("output = %s\n", strings.Join(output, "|"))
	}
	fmt.Println()
}
