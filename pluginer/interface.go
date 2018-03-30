package pluginer

// Methods interface for method stub
type Methods interface {
	PluginName() string
	DefaultMethod(input []byte) error
	GetMethodByName(name string) (func(input []byte) error, error)
}
