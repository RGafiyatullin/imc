package config

type NetConfig interface {
	BindSpec() string
}

type Config interface {
	Net() NetConfig
}

func New() Config {
	netConfig := new(netConfig)
	config := new(config)
	config.netConfig = netConfig
	return config
}

type netConfig struct{}

func (this *netConfig) BindSpec() string { return ":6379" }

type config struct {
	netConfig *netConfig
}

func (this *config) Net() NetConfig { return this.netConfig }
