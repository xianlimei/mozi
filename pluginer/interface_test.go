package pluginer

import (
	"reflect"
	"testing"
)

type ImpStruct struct{}

func (i *ImpStruct) PluginName() string {
	return "test"
}

func (i *ImpStruct) DefaultMethod(input []byte) error {
	return nil
}

func (i *ImpStruct) GetMethodByName(name string) (func(input []byte) error, error) {
	return nil, nil
}

type NotImpStruct struct{}

func TestImplementMethodInterfaceByReflect(t *testing.T) {
	methods := reflect.TypeOf((*Methods)(nil)).Elem()

	s := &ImpStruct{}
	sType := reflect.TypeOf(s)
	if ok := sType.Implements(methods); !ok {
		t.Errorf("struct %s not implement interface Methods", sType.Name())
	}

	n := &NotImpStruct{}
	nType := reflect.TypeOf(n)
	if ok := nType.Implements(methods); ok {
		t.Errorf("struct %s should not implement interface Methods", nType.Name())
	}
}

func TestImplementMethodInterfaceByComplier(t *testing.T) {
	var _ Methods = (*ImpStruct)(nil) // check implement by compiler
}
