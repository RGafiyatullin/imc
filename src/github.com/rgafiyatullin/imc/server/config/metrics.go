package config

import (
	"net"
)

type MetricsConfig interface {
	GraphiteEnabled() bool
	GraphiteAddr() *net.TCPAddr
	GraphitePrefix() string
}

type metricsConfig struct{}

func (this *metricsConfig) GraphiteEnabled() bool  { return true }
func (this *metricsConfig) GraphitePrefix() string { return "imcd" }
func (this *metricsConfig) GraphiteAddr() *net.TCPAddr {
	addr, _ := net.ResolveTCPAddr("tcp", "192.168.99.100:2003")
	return addr
}
