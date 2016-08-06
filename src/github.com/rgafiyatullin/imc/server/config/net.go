package config

import "os"

const NetBindSpecDefault = ":6379"

type NetConfig interface {
	ResetToDefaults()
	BindSpec() string
}

type netConfig struct {
	bindSpec string
}

func (this *netConfig) BindSpec() string { return this.bindSpec }
func (this *netConfig) ResetToDefaults() {
	this.bindSpec = NetBindSpecDefault
}

func (this *netConfig) ReadFromOSEnv() {
	netBind := os.Getenv("IMCD_NET_BIND")
	if netBind == "" {
		this.bindSpec = NetBindSpecDefault
	} else {
		this.bindSpec = netBind
	}
}
