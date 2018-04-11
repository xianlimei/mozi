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
	Element []int `json:"element"`
}
