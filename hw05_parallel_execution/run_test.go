package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("border cases", func(t *testing.T) {
		var runTasksCount int32
		tCases := []struct {
			name string
			N    int
			M    int
			lenT int
			err  error
		}{
			{"len(tests)<N", 10, 1, 5, nil},
			{"N<=0", 0, 1, 5, fmt.Errorf("at least one goroutine and one possible error is needed to complete tasks")},
			{"M<=0", 1, 0, 5, fmt.Errorf("at least one goroutine and one possible error is needed to complete tasks")},
			{"no tasks", 10, 1, 0, nil},
		}
		for _, tCase := range tCases {
			tasks := make([]Task, 0, tCase.lenT)
			for i := 0; i < tCase.lenT; i++ {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return fmt.Errorf("error from task %d", i)
				})
			}
			result := Run(tasks, tCase.N, tCase.M)
			require.Equalf(t, tCase.err, result, tCase.name)
		}
	})
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return fmt.Errorf("error from task %d", i)
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}
