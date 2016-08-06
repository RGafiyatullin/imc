package config

import (
	"os"
	"strconv"
)

const NetBindSpecDefault = ":6379"
const NetDefaultAcceptorsCount = 1

type NetConfig interface {
	ResetToDefaults()
	BindSpec() string
	Password() string
	AcceptorsCount() int
}

type netConfig struct {
	bindSpec       string
	password       string
	acceptorsCount int
}

func (this *netConfig) AcceptorsCount() int { return this.acceptorsCount }
func (this *netConfig) BindSpec() string    { return this.bindSpec }
func (this *netConfig) Password() string    { return this.password }

func (this *netConfig) ResetToDefaults() {
	this.bindSpec = NetBindSpecDefault
	this.acceptorsCount = NetDefaultAcceptorsCount
	this.password = ""
}

func (this *netConfig) ReadFromOSEnv() {
	netBind := os.Getenv("IMCD_NET_BIND")
	if netBind == "" {
		this.bindSpec = NetBindSpecDefault
	} else {
		this.bindSpec = netBind
	}
	acceptorsCount, err := strconv.ParseInt(os.Getenv("IMCD_NET_ACCEPTORS_COUNT"), 10, 32)
	if err != nil {
		this.acceptorsCount = NetDefaultAcceptorsCount
	} else {
		this.acceptorsCount = int(acceptorsCount)
	}

	this.password = os.Getenv("IMCD_PASSWORD")
}
