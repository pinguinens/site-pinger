package config

import "time"

type Dialer struct {
	Timeout   time.Duration `yaml:"timeout" default:"30s"`
	KeepAlive time.Duration `yaml:"keepalive" default:"30s"`
}
