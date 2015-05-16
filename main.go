package main

import (
	"github.com/telemetryapp/nginxwatch/nginx"
	"gopkg.in/alecthomas/kingpin.v1"
	"log"
	"os"
	"sync"
)

var (
	configFile = kingpin.Flag("config", "Filename for the config file").Required().Short('c').ExistingFile()
)

func main() {

	kingpin.Version("1.0.0")
	kingpin.Parse()

	hostname, hosterr := os.Hostname()
	if hosterr != nil {
		log.Fatal("Unable to determine hostname")
	}

	config := nginx.ReadConfig(*configFile)

	var wg sync.WaitGroup

	// Global Report
	go nginx.TailNginx(config.Nginx, config.Graphite, nginx.Report{}, hostname)
	wg.Add(1)

	// Specific Optional Reports
	for _, report := range config.Reports {
		go nginx.TailNginx(config.Nginx, config.Graphite, report, hostname)
		wg.Add(1)
	}

	wg.Wait()
	log.Print("All log file tailing complete, exiting")
}
