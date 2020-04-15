package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, n int, m int) (merr error) {
	if n <= 0 || m <= 0 {
		return errors.New("at least one goroutine and one possible error is needed to complete tasks")
	}
	var errCount int64
	Mvalue := int64(m)
	wg := &sync.WaitGroup{}
	goroutCount := n
	if len(tasks) < n {
		goroutCount = len(tasks) //не запускаем лишние горутины
	}
	donech := make(chan struct{})
	jobchan := make(chan Task)
	wg.Add(goroutCount)
	for i := 0; i < goroutCount; i++ {
		go func(done chan struct{}, jobs chan Task) {
			defer func() {
				wg.Done() //работает через замыкание
			}()
			for {
				select {
				case <-done:
					return //сигнал прекратить работу
				default:
					job, ok := <-jobs
					if !ok {
						return //канал закрыт
					}
					if err := job(); err != nil {
						atomic.AddInt64(&errCount, 1)
					}
				}
			}
		}(donech, jobchan)
	}
	for _, task := range tasks {
		if atomic.LoadInt64(&errCount) < Mvalue {
			jobchan <- task
			continue
		}
		donech <- struct{}{}
		merr = ErrErrorsLimitExceeded
		break
	}
	close(jobchan)
	wg.Wait()
	return
}
