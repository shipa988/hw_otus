package hw06_pipeline_execution //nolint:golint,stylecheck

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(name string, f func(v I) I) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v I) I { return v }),
		g("Multiplier (* 2)", func(v I) I { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v I) I { return v.(int) + 100 }),
		g("Stringifier", func(v I) I { return strconv.Itoa(v.(int)) }),
	}
	t.Run("signature change case", func(t *testing.T) {
		type realIn = <-chan I
		type realOut = <-chan I
		require.IsType(t, new(realIn), new(In), "We said: No changing the signature")   //входящий канал только на чтение
		require.IsType(t, new(realOut), new(Out), "We said: No changing the signature") //исходящий канал только на чтение
		out := ExecutePipeline(*new(Bi), *new(Bi), stages...)
		require.IsType(t, *new(realOut), out, "We said: No changing the signature") //возвращаемый канал только на чтение
	})
	t.Run("empty stages case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]int, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, []Stage{}...) {
			result = append(result, s.(int))
		}
		elapsed := time.Since(start)

		require.Equal(t, result, []int{1, 2, 3, 4, 5})
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})
	t.Run("no data case", func(t *testing.T) {
		in := make(Bi)
		data := []int{}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, result, []string{})
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})
	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, result, []string{"102", "104", "106", "108", "110"})
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})
	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			done <- <-time.After(abortDur)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}
