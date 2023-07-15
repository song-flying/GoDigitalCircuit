package digital

import (
	"context"
	"fmt"
	"goos/bit"
	"os"
	"strings"
	"sync"
	"time"
)

type Source struct {
	C      chan []bit.Bit
	values []bit.Bit
	inputs []chan bit.Bit
	tick   <-chan time.Time
}

func NewSource(ctx context.Context, valueSource chan []bit.Bit, inputs []chan bit.Bit) Source {
	if len(inputs) == 0 {
		panic("need at least one input")
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	source := Source{
		C:      valueSource,
		inputs: inputs,
		tick:   ticker.C,
	}

	var output []string
	go func() {
		for {
			select {
			case <-ctx.Done():
				source.Close()
				return
			case newValues, ok := <-valueSource:
				if !ok {
					source.Close()
					return
				}
				if len(newValues) != len(inputs) {
					panic(fmt.Sprintf("number of values and input channels don't match: %d != %d", len(newValues), len(inputs)))
				}
				source.values = newValues
			case _, ok := <-source.tick:
				if !ok {
					source.Close()
					return
				}
				if len(source.values) > 0 {
					output = append(
						output,
						strings.Join(bit.BitSliceToStringSlice(source.values), ""))
					fmt.Fprintf(os.Stderr, "\rvalues = %s", strings.Join(output, "|"))
					source.sendInput()
				}
			}
		}
	}()

	return source
}

func (s *Source) sendInput() {
	var wg sync.WaitGroup

	for i, v := range s.values {
		wg.Add(1)
		go func(i int, v bit.Bit) {
			defer wg.Done()
			s.inputs[i] <- v
		}(i, v)
	}

	wg.Wait()
}

func (s *Source) Close() {
	for _, c := range s.inputs {
		close(c)
	}
}
