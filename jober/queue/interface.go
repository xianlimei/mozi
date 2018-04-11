package queue

import (
	"github.com/go-done/mozi/jober/structs"
)

// Queue interface
type Queue interface {
	Add(job *structs.JobArgs) error
	Pop() (*structs.JobArgs, error)
}
