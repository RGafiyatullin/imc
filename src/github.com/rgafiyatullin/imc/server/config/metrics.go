package config

import (
	"net"
	"os"
)

type MetricsConfig interface {
	GraphiteEnabled() bool
	GraphiteAddr() *net.TCPAddr
	GraphitePrefix() string
}

type metricsConfig struct {
	graphiteEnabled bool
	graphitePrefix  string
	graphiteAddr    string
}

func (this *metricsConfig) GraphiteEnabled() bool  { return this.graphiteEnabled }
func (this *metricsConfig) GraphitePrefix() string { return this.graphitePrefix }
func (this *metricsConfig) GraphiteAddr() *net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", this.graphiteAddr)
	if err != nil {
		return nil
	} else {
		return addr
	}
}

func (this *metricsConfig) ResetToDefaults() {
	this.graphiteEnabled = false
	this.graphiteAddr = ""
	this.graphitePrefix = "imcd"
}

func (this *metricsConfig) ReadFromOSEnv() {
	this.graphiteAddr = os.Getenv("IMCD_GRAPHITE_ADDR_PLAINTEXT")
	this.graphiteEnabled = this.graphiteAddr != ""
	this.graphitePrefix = os.Getenv("IMCD_GRAPHITE_PREFIX")
}
