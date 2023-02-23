package site

type Target struct {
	Method string   `yaml:"method"`
	URI    string   `yaml:"uri"`
	Hosts  []string `yaml:"hosts"`
}
