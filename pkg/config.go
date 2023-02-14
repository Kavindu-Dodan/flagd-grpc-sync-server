package Core

import "flag"

const hostDefault = "localhost"
const portDefault = "9090"

type Config struct {
	Host     string
	Port     string
	Secure   bool
	CertPath string
	KeyPath  string
}

func LoadConfig() Config {
	cfg := Config{}

	flag.StringVar(&cfg.Host, "h", hostDefault, "hostDefault of the server")
	flag.StringVar(&cfg.Port, "p", portDefault, "portDefault of the server")
	flag.BoolVar(&cfg.Secure, "s", false, "enable tls")
	flag.StringVar(&cfg.CertPath, "certPath", "", "certificate path for tls connection")
	flag.StringVar(&cfg.KeyPath, "keyPath", "", "certificate key for tls connection")
	flag.Parse()

	return cfg
}
