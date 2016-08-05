package config

type StorageConfig interface {
	ResetToDefaults()
	RingSize() uint
}

type Config interface {
	ResetToDefaults()
	Net() NetConfig
	Storage() StorageConfig
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

type config struct {
	netConfig     *netConfig
	storageConfig *storageConfig
}

func (this *config) Net() NetConfig         { return this.netConfig }
func (this *config) Storage() StorageConfig { return this.storageConfig }
func (this *config) ResetToDefaults() {
	this.netConfig.ResetToDefaults()
	this.storageConfig.ResetToDefaults()
}
