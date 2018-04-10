package jober

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-done/mozi/jober/extracter"
	"github.com/go-done/mozi/jober/queue"
	"github.com/go-done/mozi/jober/structs"
	"github.com/go-done/mozi/notify"
	"github.com/go-done/mozi/pluginer"
	"github.com/go-done/mozi/util"
)

// Jober job manage
type Jober struct {
	jobs     map[string]*Job
	dir      string // directory for job file
	sodir    string // directory for so file
	plger    *pluginer.Pluginer
	hotLoad  bool                  // hot to load job file
	notifyer *notify.FileWatcher   // filewatch
	exter    extracter.Extracter   // extracter
	jobch    chan *structs.JobArgs // chan for job to be exec
	queue    queue.Queue           // job queue
}

// NewJober create a new Job
func NewJober(dir, sodir string) *Jober {
	plger := pluginer.NewPluginer(sodir)
	notifyer := notify.NewFileWatcher()
	notifyer.AddDir(dir)
	exter := extracter.NewExtracter()
	q := queue.NewRedisQueue("localhost:6379", "", 0)

	return &Jober{
		jobs:     make(map[string]*Job),
		dir:      dir,
		sodir:    sodir,
		plger:    plger,
		hotLoad:  true,
		notifyer: notifyer,
		exter:    exter,
		jobch:    make(chan *structs.JobArgs, 20),
		queue:    q,
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

	// load job from dir
	j.loadJobsFromDir()

	// start load job
	go j.jobLoader(j.jobch)

	// start worker
	for i := 0; i < 10; i++ {
		go j.jobWorker(i, j.jobch)
	}

	time.Sleep(30 * time.Second)
}

// jobWorker worker for run job
func (j *Jober) jobWorker(idx int, jobch <-chan *structs.JobArgs) error {
	for {
		job, ok := <-jobch

		if ok && job != nil && job.Args != nil && job.Name != "" {
			fmt.Printf("=== jobWorker [%d] === \n", idx)
			err := j.execJob(job) // TODO send back exec result to another channel
			fmt.Println("execjob error:", err)
		}
	}
}

// AddJob add a job
func (j *Jober) AddJob(jobBody []byte) error {
	jobArgs, err := j.exter.ExtractJob(jobBody)
	if err != nil {
		return err
	}
	return j.exter.SendJobToQueue(j.queue, jobArgs)
}

// jobLoader load job from queue
func (j *Jober) jobLoader(jobch chan<- *structs.JobArgs) {
	for {
		job, err := j.exter.LoadJobFromQueue(j.queue)
		time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Println("LoadJobFromQueue error: ", err)
			continue
		}
		jobch <- job
	}
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

type Payload struct {
	Element []int `json:"element"`
}

// execJob run a job
func (j *Jober) execJob(args *structs.JobArgs) error {
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
	// p := &Payload{}
	// json.Unmarshal(args.Args, p)
	// fmt.Printf("exec job, %+v", p)
	fmt.Printf("exec job, %+v, %T\n", args, args)
	//fmt.Printf("exec job, %+v, %T\n", string(args.Args), args.Args)
	// run the job
	job.Run(name, args.Args)

	return nil
}

// Clear clear all plugin
func (j *Jober) Clear() {
	j.plger.DestroyAllPlugins()
	os.RemoveAll(j.sodir)
}

func getPluginName(jobName string) string {
	return strings.Split(jobName, ".")[0]
}
