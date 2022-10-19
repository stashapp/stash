package job

import (
	"context"

	"github.com/remeh/sizedwaitgroup"
)

type taskExec struct {
	task
	fn func(ctx context.Context)
}

type TaskQueue struct {
	p     *Progress
	wg    sizedwaitgroup.SizedWaitGroup
	tasks chan taskExec
	done  chan struct{}
}

func NewTaskQueue(ctx context.Context, p *Progress, queueSize int, processes int) *TaskQueue {
	ret := &TaskQueue{
		p:     p,
		wg:    sizedwaitgroup.New(processes),
		tasks: make(chan taskExec, queueSize),
		done:  make(chan struct{}),
	}

	go ret.executer(ctx)

	return ret
}

func (tq *TaskQueue) Add(description string, fn func(ctx context.Context)) {
	tq.tasks <- taskExec{
		task: task{
			description: description,
		},
		fn: fn,
	}
}

func (tq *TaskQueue) Close() {
	close(tq.tasks)
	// wait for all tasks to finish
	<-tq.done
}

func (tq *TaskQueue) executer(ctx context.Context) {
	defer close(tq.done)
	defer tq.wg.Wait()
	for task := range tq.tasks {
		if IsCancelled(ctx) {
			return
		}

		tt := task

		tq.wg.Add()
		go func() {
			defer tq.wg.Done()
			tq.p.ExecuteTask(tt.description, func() {
				tt.fn(ctx)
			})
		}()
	}
}
