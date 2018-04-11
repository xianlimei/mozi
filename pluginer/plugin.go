package pluginer

import (
	"fmt"
	"plugin"
	"strings"

	"github.com/gobuffalo/uuid"
)

// RunFunc name for lookup
const RunFunc = "Core"

// Plugin a plugin instance
type Plugin struct {
	name    string
	uuid    string
	path    string // so 文件路径
	plg     *plugin.Plugin
	methods Methods
}

// NewPlugin create a new plugin
func NewPlugin(path string) (*Plugin, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("uuid.NewV4 error: %s", err.Error())
	}
	plugin := &Plugin{
		path: path,
		uuid: uuid.String(),
	}
	if err := plugin.init(); err != nil {
		return nil, fmt.Errorf("plugin init error: %s", err.Error())
	}
	return plugin, nil
}

// GetPath get the path of plugin's so file
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

// init get plugin
func (p *Plugin) init() error {
	plg, err := plugin.Open(p.path)
	if err != nil {
		return fmt.Errorf("could not open plugin at %s: %s", p.path, err.Error())
	}

	// bind core function
	symbol, err := plg.Lookup(RunFunc)
	if err != nil {
		return err
	}

	// check func type
	itf, ok := symbol.(func() interface{})
	if !ok {
		return fmt.Errorf("plugin %s:: function [Methods] should be <func() interface{}>", p.path)
	}
	// run and check the return value, the value should implemet Methods interface
	methods, ok := itf().(Methods)
	if !ok {
		return fmt.Errorf("plugin %s is not implement Methods interface", p.path)
	}

	// bind plugin name
	p.name = methods.PluginName()
	p.methods = methods

	// finally bind plugin
	p.plg = plg
	return nil
}

// Run run job with the loaded plugin
func (p *Plugin) Run(jobName string, input []byte) error {
	vars := strings.Split(jobName, ".")
	if vars[0] != p.GetName() {
		return fmt.Errorf("plugin %s is not found", vars[0])
	}
	switch len(vars) {
	case 1:
		return p.runDefault(input)
	case 2:
		return p.runByName(vars[1], input)
	}
	return fmt.Errorf("job %s is not a valid job", jobName)
}

// runByName run job with the loaded plugin
func (p *Plugin) runByName(name string, input []byte) error {
	method, err := p.methods.GetMethodByName(name)
	if err != nil {
		return err
	}
	return method(input)
}

// runDefault run default method
func (p *Plugin) runDefault(input []byte) error {
	return p.methods.DefaultMethod(input)
}

func (p *Plugin) String() string {
	return fmt.Sprintf("plugin id: %s, plugin name: %s, path: %s", p.GetID(), p.GetName(), p.GetPath())
}
