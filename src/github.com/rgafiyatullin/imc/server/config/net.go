package config

type NetConfig interface {
	ResetToDefaults()
	BindSpec() string
}

type netConfig struct {
	bindSpec string
}

func (this *netConfig) BindSpec() string { return this.bindSpec }
func (this *netConfig) ResetToDefaults() {
	this.bindSpec = ":6379"
}
