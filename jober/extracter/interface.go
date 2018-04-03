package extracter

import (
	"github.com/go-done/mozi/jober/queue"
	"github.com/go-done/mozi/jober/structs"
)

// Extracter interface
type Extracter interface {
	ExtractJob(body []byte) (*structs.JobArgs, error)
	SendJobToQueue(q queue.Queue, jobArgs *structs.JobArgs) error
	LoadJobFromQueue(q queue.Queue) (*structs.JobArgs, error)
}
