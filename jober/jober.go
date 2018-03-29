package jober

import (
	"fmt"
	"os"

	"github.com/chenkaiC4/golang-plugins/notify"
	"github.com/chenkaiC4/golang-plugins/pluginer"
	"github.com/chenkaiC4/golang-plugins/util"
)

// Jober job manage
type Jober struct {
	jobs     map[string]*Job
	dir      string // directory for job file
	sodir    string // directory for so file
	plger    *pluginer.Pluginer
	hotLoad  bool                // hot to load job file
	notifyer *notify.FileWatcher // filewatch
}

// NewJober create a new Job
func NewJober(dir, sodir string) *Jober {
	plger := pluginer.NewPluginer(sodir)
	notifyer := notify.NewFileWatcher()
	notifyer.AddDir(dir)

	return &Jober{
		jobs:     make(map[string]*Job),
		dir:      dir,
		sodir:    sodir,
		plger:    plger,
		hotLoad:  true,
		notifyer: notifyer,
	}
}

// Start begin process
func (j *Jober) Start() {
	if j.hotLoad {
		go func() {
			for signal := range j.notifyer.ChangedFile {
				_, fp, err := notify.DecodeSignal(signal)
				if err != nil {
					continue
				}
				j.plger.LoadPlugin(fp)
			}
		}()
	}

	j.loadJobsFromDir()
}

//ExecPluginMethodByName exec plugin method by name
func (j *Jober) ExecPluginMethodByName(name string) {
	j.plger.RunMethodByName(name)
}

func (j *Jober) loadJobsFromDir() (err error) {
	util.TraversalDir(j.dir, func(fp string) {
		fmt.Printf("begin to load job: %s\nresult: \t", fp)
		if err == nil {
			e := j.plger.LoadPlugin(fp)
			if e != nil {
				fmt.Println("Failed")
				err = fmt.Errorf("load job: %s failed, err: %v", fp, e)
				return
			}
			fmt.Println("Success")
		}
	})
	return err
}

// RunJobByID run a job by ID
func (j *Jober) RunJobByID(id string) {
	j.jobs[id].Run()
}

// AddJob add a new Job
func (j *Jober) AddJob(job *Job) error {
	if _, has := j.jobs[job.ID]; !has {
		j.jobs[job.ID] = job
	}
	return nil
}

// Clear clear all plugin
func (j *Jober) Clear() {
	j.plger.DestroyAllPlugins()
	os.RemoveAll(j.sodir)
}
