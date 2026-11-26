package job

import (
	"context"
	"sync"
)

type taskExec struct {
	task
	fn func(ctx context.Context)
}

type TaskQueue struct {
	p     *Progress
	wg    sync.WaitGroup
	slots chan struct{}
	tasks chan taskExec
	done  chan struct{}
}

// NewTaskQueue returns a queue that runs its tasks concurrently. A task runs
// only while it holds a slot in slots, so queues that share slots never run
// more than cap(slots) tasks at once between them.
func NewTaskQueue(ctx context.Context, p *Progress, queueSize int, slots chan struct{}) *TaskQueue {
	ret := &TaskQueue{
		p:     p,
		slots: slots,
		tasks: make(chan taskExec, queueSize),
		done:  make(chan struct{}),
	}

	go ret.executer(ctx)

	return ret
}

func (tq *TaskQueue) Add(description string, fn func(ctx context.Context)) {
	tq.p.AddTotal(1)
	tq.tasks <- taskExec{
		task: task{
			description: description,
		},
		fn: fn,
	}
}

func (tq *TaskQueue) Close() {
	close(tq.tasks)
	tq.p.Definite()
}

func (tq *TaskQueue) Wait() {
	// wait for all tasks to finish
	<-tq.done
}

func (tq *TaskQueue) executer(ctx context.Context) {
	defer close(tq.done)
	defer tq.wg.Wait()
	for task := range tq.tasks {
		if IsCancelled(ctx) {
			continue // allow channel to continue draining until Close()
		}

		tt := task

		select {
		case tq.slots <- struct{}{}:
		case <-ctx.Done():
			continue
		}

		tq.wg.Add(1)
		go func() {
			defer func() { <-tq.slots }()
			defer tq.wg.Done()
			defer tq.p.Increment()
			tq.p.ExecuteTask(tt.description, func() {
				tt.fn(ctx)
			})
		}()
	}
}
