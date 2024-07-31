package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hibiken/asynq"
	"github.com/oxiginedev/gostarter/internal/pkg/redis"
	"github.com/oxiginedev/gostarter/internal/worker/task"
	"github.com/oxiginedev/gostarter/util"
	"github.com/pkg/errors"
)

const (
	QueueHigh    = "high"
	QueueDefault = "default"
	QueueLow     = "low"
)

var (
	ErrQueueNotFound = fmt.Errorf("asynq: %w", asynq.ErrQueueNotFound)
	ErrTaskNotFound  = fmt.Errorf("asynq: %w", asynq.ErrTaskNotFound)
)

type Job struct {
	// ID is the unique identifier of the job
	ID uuid.UUID `json:"id"`
	// Queue the job should be pushed on to
	Queue string `json:"queue"`
	// Args that will be passed to the Handler when called
	Payload interface{}   `json:"payload"`
	Delay   time.Duration `json:"delay"`
}

type QueueConfig struct {
	// Queues is a map of available queues with corresponding priority
	Queues map[string]int
	// Type
	Type         string
	RedisClient  *redis.Redis
	RedisAddress string
}

type Queue struct {
	client    *asynq.Client
	config    *QueueConfig
	inspector *asynq.Inspector
}

func NewQueue(cfg *QueueConfig) *Queue {
	client := asynq.NewClient(cfg.RedisClient)
	inspector := asynq.NewInspector(cfg.RedisClient)

	return &Queue{
		client:    client,
		config:    cfg,
		inspector: inspector,
	}
}

func (q *Queue) Enqueue(taskName task.TaskName, job *Job) error {
	if util.IsStringEmpty(job.ID.String()) {
		job.ID = uuid.Must(uuid.NewV7())
	}

	if util.IsStringEmpty(job.Queue) {
		job.Queue = QueueDefault
	}

	payload, err := json.Marshal(job.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal task args")
	}

	task := asynq.NewTask(taskName.String(), payload,
		asynq.Queue(job.Queue),
		asynq.TaskID(job.ID.String()),
		asynq.ProcessIn(job.Delay),
	)

	_, err = q.inspector.GetTaskInfo(job.Queue, job.ID.String())
	if err != nil {
		message := err.Error()
		if ErrQueueNotFound.Error() == message || ErrTaskNotFound.Error() == message {
			_, err := q.client.Enqueue(task, nil)
			return err
		}

		return err
	}

	err = q.inspector.DeleteTask(job.Queue, job.ID.String())
	if err != nil {
		return err
	}

	_, err = q.client.Enqueue(task, nil)
	return err
}
