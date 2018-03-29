package pluginer

import (
	"sync"

	"github.com/chenkaiC4/golang-plugins/util"
)

// Pluginer use golang plugin package to hot load program
type Pluginer struct {
	plgs map[string]*Plugin // string is id of plugin
	p2id map[string]string  // k-v: <path, plugin id>
	dir  string
	mux  *sync.Mutex
}

// NewPluginer create pluginer
func NewPluginer(dir string) *Pluginer {
	util.MustBeDir(dir)
	return &Pluginer{
		plgs: make(map[string]*Plugin),
		p2id: make(map[string]string),
		dir:  dir,
		mux:  new(sync.Mutex),
	}
}

// LoadPlugin load a new plugin by path
func (p *Pluginer) LoadPlugin(path string) error {
	plg, err := NewPlugin(path, p.dir)
	if err != nil {
		return err
	}
	return p.AddPlugin(plg)
}

// AddPlugin add a plugin
func (p *Pluginer) AddPlugin(plg *Plugin) error {

	id := plg.GetID()
	path := plg.GetPath()

	// bind kv for lookup by path
	if _, exist := p.p2id[path]; !exist {
		// path not exist, a new plugin
		p.p2id[path] = id
		p.plgs[id] = plg
	} else {
		// path exist already, update
		oldPluginID := p.p2id[path]
		p.plgs[oldPluginID].Destroy() // remove old plugin
		delete(p.plgs, oldPluginID)   // delete old id
		p.plgs[id] = plg              // add new id
		p.p2id[path] = id             // update path id
	}

	return nil
}

// RunMethodByName run method by name
func (p *Pluginer) RunMethodByName(name string) {
	for k := range p.plgs {
		p.plgs[k].RunMethodByName(name)
	}

	return
}

// DestroyAllPlugins destroy all plugins
func (p *Pluginer) DestroyAllPlugins() error {
	for i := range p.plgs {
		if err := p.plgs[i].Destroy(); err != nil {
			return err
		}
	}
	return nil
}
