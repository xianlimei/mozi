package jober

import (
	"github.com/go-done/mozi/pluginer"
	"github.com/satori/go.uuid"
)

// Job a job
type Job struct {
	id     string
	plugin *pluginer.Plugin
}

// NewJob create a new job
func NewJob(plg *pluginer.Plugin) *Job {
	return &Job{
		id:     uuid.NewV4().String(),
		plugin: plg,
	}
}

// GetID get job id
func (j *Job) GetID() string {
	return j.id
}

// Run job
func (j *Job) Run(jobName string, input []byte) error {
	return j.plugin.Run(jobName, input)
}
