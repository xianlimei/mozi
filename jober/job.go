package jober

// Job a job
type Job struct {
	ID    string
	Name  string
	execs []func()
}

// Run job
func (j *Job) Run() {
	for i := range j.execs {
		j.execs[i]()
	}
}
