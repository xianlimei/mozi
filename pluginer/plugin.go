package pluginer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/satori/go.uuid"
)

// RunFunc name for lookup
const RunFunc = "Core"

// Plugin a plugin program
type Plugin struct {
	name    string
	uuid    string
	path    string // 源文件路径
	dst     string // so 文件夹路径
	sopath  string // so path
	plg     *plugin.Plugin
	methods Methods
}

// NewPlugin create a new program
func NewPlugin(path, dst string) (*Plugin, error) {
	program := &Plugin{
		path: path,
		dst:  dst,
	}
	if err := program.init(); err != nil {
		return nil, err
	}
	return program, nil
}

// GetPath get program path
func (p *Plugin) GetPath() string {
	return p.path
}

// GetID get id
func (p *Plugin) GetID() string {
	return p.uuid
}

// GetName get name
func (p *Plugin) GetName() string {
	return p.name
}

// GetDst get dst
func (p *Plugin) GetDst() string {
	return p.dst
}

// init get plugin
func (p *Plugin) init() error {
	p.uuid = uuid.NewV4().String()
	p.sopath = filepath.Join(p.dst, p.uuid) + ".so"

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o="+p.sopath, p.path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not compile %s: %v", p.path, err)
	}

	plg, err := plugin.Open(p.sopath)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", p.sopath, err)
	}

	// bind core function
	symbol, err := plg.Lookup(RunFunc)
	if err != nil {
		return err
	}

	// check func type
	itf, ok := symbol.(func() interface{})
	if !ok {
		return fmt.Errorf("plugin %s:: function [Methods] should be <func() interface{}>", p.sopath)
	}
	// run and check the return value, the value should implemet Methods interface
	methods, ok := itf().(Methods)
	if !ok {
		return fmt.Errorf("plugin %s is not implement Methods interface", p.sopath)
	}
	p.name = methods.PluginName()
	p.methods = methods
	// finally bind plugin
	p.plg = plg

	return nil
}

// Run run job with the loaded plugin
func (p *Plugin) Run(jobName string, input []byte) error {
	vars := strings.Split(jobName, ".")
	switch len(vars) {
	case 1:
		return p.RunDefault(input)
	case 2:
		return p.RunByName(vars[1], input)
	}
	return p.methods.DefaultMethod(input)
}

// RunByName run job with the loaded plugin
func (p *Plugin) RunByName(name string, input []byte) error {
	method, err := p.methods.GetMethodByName(name)
	if err != nil {
		return err
	}
	return method(input)
}

// RunDefault run default method
func (p *Plugin) RunDefault(input []byte) error {
	fmt.Println("RunDefault 执行了")
	return p.methods.DefaultMethod(input)
}

func (p *Plugin) String() string {
	return fmt.Sprintf("plugin %s [%s], path: %s", p.name, p.uuid, p.path)
}

// Destroy destroy plugin
func (p *Plugin) Destroy() error {
	p.plg = nil
	fp := filepath.Join(".", p.sopath)
	err := os.Remove(fp)
	return err
}
