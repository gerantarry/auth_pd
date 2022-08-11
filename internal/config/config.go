package config

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		BindIp string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"listen"`
	Database struct {
		Scheme   string `yaml:"scheme"`
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
		BindIp   string `yaml:"bind_ip"`
		port     string `yaml:"port"`
	}
}
