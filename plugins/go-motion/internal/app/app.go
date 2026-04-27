package app

type Config struct {
	PluginName string
}

func NewDefaultConfig() Config {
	return Config{
		PluginName: "go-motion",
	}
}
