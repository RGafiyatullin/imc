package config

type NetConfig interface {
	ResetToDefaults()
	BindSpec() string
}

type StorageConfig interface {
	ResetToDefaults()
	RingSize() uint
}

type Config interface {
	ResetToDefaults()
	Net() NetConfig
}

func New() Config {
	netConfig := new(netConfig)
	storageConfig := new(storageConfig)
	config := new(config)
	config.netConfig = netConfig
	config.storageConfig = storageConfig
	config.ResetToDefaults()
	return config
}

type netConfig struct{
	bindSpec string
}
func (this *netConfig) BindSpec() string { return this.bindSpec }
func (this *netConfig) ResetToDefaults() {
	this.bindSpec = ":6379"
}


type storageConfig struct {
	ringSize uint
}
func (this *storageConfig) RingSize() uint { return this.ringSize }
func (this *storageConfig) ResetToDefaults() {
	this.ringSize = 32
}

type config struct {
	netConfig *netConfig
	storageConfig *storageConfig
}

func (this *config) Net() NetConfig { return this.netConfig }
func (this *config) Storage() StorageConfig { return this.storageConfig }
func (this *config) ResetToDefaults() {
	this.netConfig.ResetToDefaults()
	this.storageConfig.ResetToDefaults()
}
