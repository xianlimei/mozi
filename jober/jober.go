package jober

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chenkaiC4/golang-plugins/notify"
	"github.com/chenkaiC4/golang-plugins/pluginer"
	"github.com/chenkaiC4/golang-plugins/util"
)

// JobArgs create for
type JobArgs struct {
	Name string `json:"name"`
	Args []byte `json:"args"`
}

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

func (j *Jober) loadJobsFromDir() (err error) {
	util.TraversalDir(j.dir, func(fp string) {
		fmt.Printf("begin to load job: %s\nresult: \t", fp)
		if err == nil {
			e := j.plger.LoadPlugin(fp)
			if e != nil {
				fmt.Printf("Failed\n %v", e)
				err = fmt.Errorf("load job: %s failed, err: %v", fp, e)
				return
			}
			fmt.Println("Success")
		}
	})
	return err
}

// AddJob add a new Job
func (j *Jober) AddJob(args *JobArgs) error {
	name := args.Name
	if name == "" {
		return errors.New("job name is empty")
	}
	plgName := getPluginName(name)
	plg, err := j.plger.GetPluginByName(plgName)
	if err != nil {
		return fmt.Errorf("find plugin by name failed: %v", err)
	}

	job := NewJob(plg)
	id := job.GetID()
	j.jobs[id] = job

	// run the job
	job.Run(name, args.Args)

	return nil
}

func getPluginName(jobName string) string {
	return strings.Split(jobName, ".")[0]
}

// Clear clear all plugin
func (j *Jober) Clear() {
	j.plger.DestroyAllPlugins()
	os.RemoveAll(j.sodir)
}
