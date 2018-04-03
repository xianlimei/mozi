package extracter

import (
	"encoding/json"

	"github.com/go-done/mozi/jober/queue"
	"github.com/go-done/mozi/jober/structs"
)

type extracter struct {
}

// Extract extract job from body, which maybe come from http, websocket, gRPC etc.
func (e *extracter) ExtractJob(body []byte) (*structs.JobArgs, error) {
	var jobArgs structs.JobArgs
	err := json.Unmarshal(body, &jobArgs)
	if err != nil {
		return nil, err
	}
	return &jobArgs, nil
}

// SendToQueue send job to queue
func (e *extracter) SendJobToQueue(q queue.Queue, jobArgs *structs.JobArgs) error {
	return q.Add(jobArgs)
}

func (e *extracter) LoadJobFromQueue(q queue.Queue) (*structs.JobArgs, error) {
	return q.Pop()
}

// NewExtracter create a extracter
func NewExtracter() Extracter {
	return &extracter{}
}
