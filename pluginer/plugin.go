package pluginer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/satori/go.uuid"
)

// Plugin a plugin program
type Plugin struct {
	id     string
	path   string // 源文件路径
	dst    string // so 文件夹路径
	sopath string // so path
	plg    *plugin.Plugin
}

// NewPlugin create a new program
func NewPlugin(path, dst string) (*Plugin, error) {
	id := uuid.NewV4().String()
	program := &Plugin{
		id:   id,
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
	return p.id
}

// GetDst get dst
func (p *Plugin) GetDst() string {
	return p.dst
}

// init get plugin
func (p *Plugin) init() error {
	p.sopath = filepath.Join(p.dst, p.id) + ".so"
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
	p.plg = plg
	return nil
}

// RunMethodByName run method by name
func (p *Plugin) RunMethodByName(name string) error {
	f, err := p.plg.Lookup(name)
	if err != nil {
		return err
	}
	f.(func() error)()
	return nil
}

func (p *Plugin) String() string {
	return fmt.Sprintf("plugin ID: %s, path: %s", p.id, p.path)
}

// Destroy destroy plugin
func (p *Plugin) Destroy() error {
	p.plg = nil
	fp := filepath.Join(".", p.sopath)
	err := os.Remove(fp)
	return err
}
