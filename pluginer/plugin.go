package pluginer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/satori/go.uuid"
)

// PluginName name for lookup
const PluginName = "PluginName"

// RunFunc name for lookup
const RunFunc = "Run"

// Plugin a plugin program
type Plugin struct {
	name   string
	uuid   string
	run    func([]byte) error
	path   string // 源文件路径
	dst    string // so 文件夹路径
	sopath string // so path
	plg    *plugin.Plugin
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

	// bind name
	name, err := plg.Lookup(PluginName)
	if err != nil {
		return err
	}
	p.name = name.(func() string)()

	// bind run function
	runFunc, err := plg.Lookup(RunFunc)
	if err != nil {
		return err
	}
	p.run = runFunc.(func([]byte) error)

	// finally bind plugin
	p.plg = plg

	return nil
}

// Run run job with the loaded plugin
func (p *Plugin) Run(input []byte) error {
	return p.run(input)
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
