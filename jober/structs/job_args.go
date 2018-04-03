package structs

// JobArgs get from client request, contain the job information.
type JobArgs struct {
	Name string `json:"name"` // name is the type of job
	Args []byte `json:"args"` // args is the json encoded input of job runner defined by specific plugin
}
