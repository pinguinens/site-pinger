package config

type Messenger struct {
	Enabled    bool     `yaml:"enabled"`
	AlarmCodes []string `yaml:"alarm_codes" default:"*"`
	Address    string   `yaml:"address" default:"localhost:8081"`
}
