package main

import (
	"log"
	"os"
	"sync"

	"github.com/telemetryapp/nginxwatch/nginx"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	configFile = kingpin.Flag("config", "Filename for the config file").Required().Short('c').ExistingFile()
	debug      = kingpin.Flag("debug", "Enable debug mode, prints to std out").Short('d').Bool()
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
	go nginx.TailNginx(config.Nginx, config.Graphite, nginx.Report{}, hostname, *debug)
	wg.Add(1)

	// Specific Optional Reports
	for _, report := range config.Reports {
		go nginx.TailNginx(config.Nginx, config.Graphite, report, hostname, *debug)
		wg.Add(1)
	}

	wg.Wait()
	log.Print("All log file tailing complete, exiting")
}
