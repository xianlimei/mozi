package pluginer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const examplePluginMath = `
package main

import (
	"encoding/json"
	"fmt"
)

type math struct {
	method map[string]func([]byte) error
}

func (m *math) PluginName() string {
	return "Math"
}

func (m *math) GetMethodByName(name string) (func([]byte) error, error) {
	f, exist := m.method[name]
	if exist {
		return f, nil
	}
	return nil, fmt.Errorf("func not exist: %s", name)
}

func (m *math) DefaultMethod(input []byte) error {
	if input == nil {
		return nil
	}

	p := &Payload{}
	err := json.Unmarshal(input, p)
	if err != nil {
		fmt.Println("err: ", err)
		return err
	}

	fmt.Printf("Default Method Exec: %+v\n", p.Element)
	return nil
}

func add(input []byte) error {
	p := new(Payload)
	json.Unmarshal(input, p)
	var sum int
	for _, v := range p.Element {
		sum += v
	}
	fmt.Printf("sum is: %d\n", sum)
	return nil
}

func multiply(input []byte) error {
	p := new(Payload)
	json.Unmarshal(input, p)
	var res = 1
	for _, v := range p.Element {
		res *= v
	}
	fmt.Printf("multiply is: %d\n", res)
	return nil
}

// Core get a stub method
func Core() interface{} {
	method := make(map[string]func([]byte) error)
	method["Add"] = add
	method["Multiply"] = multiply
	return &math{
		method: method,
	}
}

func main() {}

// Payload struct for input slice
type Payload struct {
	Element []int` + " `" + `json:"element"` + "`" + `
}` // tricks for back quote can't include back qoute

// Payload just for test
type Payload struct {
	Element []int `json:"element"`
}

func TestPlugin(t *testing.T) {
	// test plugin.go

	os.Mkdir("test", os.ModePerm)
	fp := filepath.Join("test", "math.go")
	fpp := filepath.Join("test", "samemath.go")
	sopath := filepath.Join("test", "math.so")
	samesopath := filepath.Join("test", "samemath.so")
	defer os.RemoveAll("test")

	fi, _ := os.Create(fp)
	fii, _ := os.Create(fpp)
	_, err := fi.WriteString(examplePluginMath)
	_, err = fii.WriteString(examplePluginMath + `
		const xxx = "123"
	`)
	if err != nil {
		t.Errorf("create math.go file failed: %s", err.Error())
	}
	fi.Close()
	fii.Close()

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o="+sopath, fp)
	if err := cmd.Run(); err != nil {
		t.Errorf("compile error")
	}

	cmd = exec.Command("go", "build", "-buildmode=plugin", "-o="+samesopath, fpp)
	if err := cmd.Run(); err != nil {
		t.Errorf("compile error")
	}

	plg, err := NewPlugin(sopath)
	if err != nil {
		t.Errorf("NewPlugin Error")
	}

	if plg.String() != fmt.Sprintf("plugin id: %s, plugin name: %s, path: %s", plg.GetID(), plg.GetName(), plg.GetPath()) {
		t.Errorf("plg String() not correct, maybe some init logic failed")
	}

	n := []int{1, 2, 3, 4, 5}
	pld := &Payload{
		Element: n,
	}
	plgb, err := json.Marshal(pld)
	if err != nil {
		t.Errorf("json payload failed: %s", err.Error())
	}
	if err := plg.Run("Math", plgb); err != nil {
		t.Errorf("run Math job failed: %s", err.Error())
	}
	if err := plg.Run("Math.Add", plgb); err != nil {
		t.Errorf("run Math.Add job failed: %s", err.Error())
	}
	if err := plg.Run("Math.Multiply", plgb); err != nil {
		t.Errorf("run Math.Multiply job failed: %s", err.Error())
	}
	if err := plg.Run("Math.NotExist", plgb); err == nil {
		t.Error("run Math.NotExist job should cause error but not")
	}
	if err := plg.Run("Mathhhh.NotExist", plgb); err == nil {
		t.Error("run Mathhhh.NotExist job should cause error but not")
	}
	if err := plg.Run("Math.Multiply.ooo", plgb); err == nil {
		t.Error("run Mathhhh.NotExist job should cause error but not")
	}

	// test pluginer.go
	plger := NewPluginer()
	if err := plger.LoadPlugin(os.TempDir() + "xxx.so"); err == nil {
		t.Errorf("pluginer load plugin at wrong path should return err")
	}
	if err := plger.LoadPlugin(sopath); err != nil {
		t.Errorf("pluginer load plugin should not return err: %s", err.Error())
	}
	_, err = plger.GetPluginByName("Mathxx")
	if err == nil {
		t.Errorf("GetPluginByName with a not exist plugin name should return err")
	}
	pllg, err := plger.GetPluginByName("Math")
	if err != nil {
		t.Errorf("GetPluginByName should not return err for Math plugin is valid")
	}
	if pllg.GetName() != plg.GetName() || pllg.GetPath() != plg.GetPath() {
		t.Errorf("should equal by name, path")
	}
	err = plger.AddPlugin(plg)
	if err != nil {
		t.Errorf("AddPlugin should not return err but with %s", err.Error())
	}

	pllg, _ = plger.GetPluginByName("Math")
	if pllg.GetID() != plg.GetID() {
		t.Errorf("update plugin, with the same path")
	}
	if err := plger.LoadPlugin(samesopath); err == nil {
		t.Errorf("load a plugin with the same name should return err")
	}
}
