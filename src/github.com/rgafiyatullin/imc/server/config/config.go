package config

type Config interface {
	ResetToDefaults()
	Net() NetConfig
	Storage() StorageConfig
	Metrics() MetricsConfig
}

func New() Config {
	netConfig := new(netConfig)
	storageConfig := new(storageConfig)
	config := new(config)
	metricsConfig := new(metricsConfig)
	config.netConfig = netConfig
	config.storageConfig = storageConfig
	config.metricsConfig = metricsConfig
	config.ResetToDefaults()
	return config
}

type config struct {
	netConfig     *netConfig
	storageConfig *storageConfig
	metricsConfig *metricsConfig
}

func (this *config) Net() NetConfig         { return this.netConfig }
func (this *config) Storage() StorageConfig { return this.storageConfig }
func (this *config) Metrics() MetricsConfig { return this.metricsConfig }
func (this *config) ResetToDefaults() {
	this.netConfig.ResetToDefaults()
	this.storageConfig.ResetToDefaults()
}
