package nginx

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Graphite Structure
type Graphite struct {
	Server   string
	Interval int64
}

// Nginx Structure
type Nginx struct {
	Logfile string
}

// Report Structure
type Report struct {
	Label    string
	Upstream string
	Host     string
	Methods  []string
	Statuses []int64
	UriRegex string
}

// Config Structure
type Config struct {
	Graphite Graphite
	Nginx    Nginx
	Reports  []Report `toml:"report"`
}

// ReadConfig reads a config file
func ReadConfig(configFile string) (config Config) {

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Panicf("Unable to parse config file: %v", err)
	}

	if config.Graphite.Interval == 0 {
		config.Graphite.Interval = 10
	}

	return config
}
