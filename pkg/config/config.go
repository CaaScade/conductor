package config

import (
	"github.com/spf13/viper"
)

var (
	Addr       = "addr"
	Port       = "port"
	KubeConfig = "kubeconfig"
)

type Config struct {
	Addr       string `json:"addr,omitempty"`
	Port       int    `json:"port,omitempty"`
	KubeConfig []byte `json:"kubeConfig,omitempty"`
}

func (c *Config) Refresh() *Config {
	c.Addr = viper.GetString(Addr)
	c.Port = viper.GetInt(Port)
	c.KubeConfig = []byte(viper.GetString(KubeConfig))
	return c
}

func GlobalConfig() *Config {
	c := &Config{}
	return c.Refresh()
}
