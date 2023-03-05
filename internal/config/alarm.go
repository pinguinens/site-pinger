package config

type Alarm struct {
	Enabled bool   `yaml:"enabled"`
	Url     string `yaml:"url" default:"localhost:8080/alarm"`
}
