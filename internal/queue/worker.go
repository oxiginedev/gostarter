package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/oxiginedev/gostarter/internal/worker/task"
)

type Worker struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

func NewWorker(ctx context.Context, concurrency int, q *Queue) *Worker {
	var opts asynq.RedisConnOpt

	srv := asynq.NewServer(
		opts,
		asynq.Config{
			Concurrency: concurrency,
			BaseContext: func() context.Context {
				return ctx
			},
			Queues: map[string]int{
				QueueHigh:    10,
				QueueDefault: 5,
				QueueLow:     3,
			},
			IsFailure: func(err error) bool {
				if _, ok := err.(*task.RateLimitError); ok {
					return false
				}
				return true
			},
			RetryDelayFunc: task.GetRetryDelay(),
		},
	)

	mux := asynq.NewServeMux()

	return &Worker{
		srv: srv,
		mux: mux,
	}
}

func (p *Worker) Start() error {
	return p.srv.Start(p.mux)
}

func (p *Worker) Register(task string, handlerFn func(context.Context, *asynq.Task) error) {
	p.mux.HandleFunc(task, asynq.HandlerFunc(handlerFn).ProcessTask)
}

func (p *Worker) Stop() {
	p.srv.Stop()
	p.srv.Shutdown()
}
