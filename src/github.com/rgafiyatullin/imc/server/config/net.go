package config

import "os"

const NetBindSpecDefault = ":6379"

type NetConfig interface {
	ResetToDefaults()
	BindSpec() string
	Password() string
}

type netConfig struct {
	bindSpec string
	password string
}

func (this *netConfig) BindSpec() string { return this.bindSpec }
func (this *netConfig) Password() string { return this.password }

func (this *netConfig) ResetToDefaults() {
	this.bindSpec = NetBindSpecDefault
	this.password = ""
}

func (this *netConfig) ReadFromOSEnv() {
	netBind := os.Getenv("IMCD_NET_BIND")
	if netBind == "" {
		this.bindSpec = NetBindSpecDefault
	} else {
		this.bindSpec = netBind
	}
	this.password = os.Getenv("IMCD_PASSWORD")
}
