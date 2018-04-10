package pluginer

import (
	"errors"
	"fmt"
)

// Pluginer use golang plugin package to hot load program
type Pluginer struct {
	path2name map[string]string // k-v: <path, name>
	name2plg  map[string]*Plugin
}

// NewPluginer create pluginer
func NewPluginer() *Pluginer {
	return &Pluginer{
		name2plg:  make(map[string]*Plugin),
		path2name: make(map[string]string),
	}
}

// LoadPlugin load a new plugin by path
func (p *Pluginer) LoadPlugin(path string) error {
	plg, err := NewPlugin(path)
	if err != nil {
		return err
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
			p := p.name2plg[name].GetPath()
			return fmt.Errorf("plugin name: [%s] is conflict with loaded plugin at %s", name, p)
		}
	}

	return nil
}

// GetPluginByName get plugin by name
func (p *Pluginer) GetPluginByName(name string) (*Plugin, error) {
	plg, exist := p.name2plg[name]
	if exist {
		return plg, nil
	}
	return nil, errors.New("plugin not found")
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
	p.name2plg[name] = plg
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
