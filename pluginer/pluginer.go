package pluginer

import (
	"errors"
	"fmt"

	"github.com/go-done/mozi/util"
)

// Pluginer use golang plugin package to hot load program
type Pluginer struct {
	path2name map[string]string // k-v: <path, name>
	name2plg  map[string]*Plugin
	dir       string
}

// NewPluginer create pluginer
func NewPluginer(dir string) *Pluginer {
	util.MustBeDir(dir)

	return &Pluginer{
		name2plg:  make(map[string]*Plugin),
		path2name: make(map[string]string),
		dir:       dir,
	}
}

// LoadPlugin load a new plugin by path
func (p *Pluginer) LoadPlugin(path string) error {
	plg, err := NewPlugin(path, p.dir)
	if err != nil {
		return fmt.Errorf("create plugin failed: %v", err)
	}
	return p.AddPlugin(plg)
}

// AddPlugin add a plugin
func (p *Pluginer) AddPlugin(plg *Plugin) error {
	path := plg.GetPath()
	name := plg.GetName()

	// name not exist  --> a new plugin
	if !p.isPluginNameExist(name) {
		p.addNewPlugin(plg)
	} else {
		// name exist but with the same path  --> update plugin
		if p.isPluginPathExist(path) {
			p.updatePluginByName(name, plg)
		} else {
			// name exist but with another path --> conflict error
			return fmt.Errorf("plugin name: [%s] is occupied at %s", name, path)
		}
	}

	return nil
}

// addNewPlugin add a pure new plugin
func (p *Pluginer) addNewPlugin(plg *Plugin) {
	name := plg.GetName()
	path := plg.GetPath()
	p.name2plg[name] = plg
	p.path2name[path] = name
}

// updatePluginByName update the plugin by name
func (p *Pluginer) updatePluginByName(name string, plg *Plugin) {
	oldPlugin := p.name2plg[name]
	p.name2plg[name] = plg
	oldPlugin.Destroy()
}

// isPluginPathExist check the plugin path existance at pluginer
func (p *Pluginer) isPluginPathExist(path string) bool {
	_, exist := p.path2name[path]
	return exist
}

// isPluginNameExist check the plugin name existance at pluginer
func (p *Pluginer) isPluginNameExist(name string) bool {
	_, exist := p.name2plg[name]
	return exist
}

// GetPluginByName get plugin by name
func (p *Pluginer) GetPluginByName(name string) (*Plugin, error) {
	plg, exist := p.name2plg[name]
	if exist {
		return plg, nil
	}
	return nil, errors.New("not found")
}

// DestroyAllPlugins destroy all plugins
func (p *Pluginer) DestroyAllPlugins() error {
	for name := range p.name2plg {
		if err := p.name2plg[name].Destroy(); err != nil {
			return err
		}
	}
	// bind new address, free the old
	p.name2plg = make(map[string]*Plugin)
	p.path2name = make(map[string]string)
	return nil
}
