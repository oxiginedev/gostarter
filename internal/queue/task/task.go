package task

import (
	"time"

	"github.com/hibiken/asynq"
)

const TaskEmailDelivery TaskName = "task:email-delivery"

type TaskName string

func (t TaskName) String() string {
	return string(t)
}

type RateLimitError struct {
	delay time.Duration
	Err   error
}

func (e *RateLimitError) Delay() time.Duration {
	return e.delay
}

func (e *RateLimitError) Error() string {
	return e.Err.Error()
}

func GetRetryDelay() asynq.RetryDelayFunc {
	return func(n int, e error, t *asynq.Task) time.Duration {
		if rateLimitErr, ok := e.(*RateLimitError); ok {
			return rateLimitErr.Delay()
		}

		return asynq.DefaultRetryDelayFunc(n, e, t)
	}
}
