package nginx

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Graphite struct {
	Server   string
	Interval int64
}

type Nginx struct {
	Logfile string
}

type Report struct {
	Label    string
	Upstream string
	Host     string
	Methods  []string
	Statuses []int64
	UriRegex string
}

type Config struct {
	Graphite Graphite
	Nginx    Nginx
	Reports  []Report `toml:"report"`
}

func ReadConfig(configFile string) (config Config) {

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Panicf("Unable to parse config file: %v", err)
	}

	return config
}
