package nginx

import (
	"fmt"
	"github.com/ActiveState/tail"
	"log"
	"net"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"time"
)

var counter uint64 = 0

var (
	statusRegex   = regexp.MustCompile(`\s+status:([^\s\\]*)`)
	timeRegex     = regexp.MustCompile(`\s+time:([^\s\\]*)`)
	sslRegex      = regexp.MustCompile(`\s+ssl:([^\s\\]*)`)
	methodRegex   = regexp.MustCompile(`\s+method:([^\s\\]*)`)
	hostRegex     = regexp.MustCompile(`\s+host:([^\s\\]*)`)
	upstreamRegex = regexp.MustCompile(`\s+upstream:([^\s\\]*)`)
	uriRegex      = regexp.MustCompile(`\s+uri:([^\s\\]*)`)
)

type nginxLogLine struct {
	time     float64
	status   int64
	ssl      string
	method   string
	host     string
	upstream string
	uri      string
}

type nginxData struct {
	hits          int64
	countNormal   int64
	countWarn     int64
	countError    int64
	percentNormal int64
	percentWarn   int64
	percentError  int64
	totalTime     float64
	minTime       float64
	maxTime       float64
	avgTime       float64
	nineTime      float64
	times         []float64
}

func (data *nginxData) reset() {
	data.hits = 0
	data.countNormal = 0
	data.countWarn = 0
	data.countError = 0
	data.percentError = 0
	data.percentNormal = 0
	data.percentWarn = 0
	data.minTime = 0
	data.maxTime = 0
	data.avgTime = 0
	data.nineTime = 0
	data.totalTime = 0
	data.times = []float64{}
}

func (data *nginxData) registerHit(line nginxLogLine) {
	data.hits++
	if line.status >= 500 {
		data.countError++
	} else if line.status >= 400 {
		data.countWarn++
	} else if line.status > 0 {
		data.countNormal++
	}

	data.percentError = int64(float64(data.countError) / float64(data.hits) * 100.0)
	data.percentNormal = int64(float64(data.countNormal) / float64(data.hits) * 100.0)
	data.percentWarn = int64(float64(data.countWarn) / float64(data.hits) * 100.0)

	if data.minTime == 0 || data.minTime > line.time {
		data.minTime = line.time
	}
	if data.maxTime < line.time {
		data.maxTime = line.time
	}
	data.totalTime += line.time
	data.avgTime = data.totalTime / float64(data.countNormal)
	data.times = append(data.times, line.time)
	tts := data.times[:]
	sort.Float64s(tts)
	index := int(float64(len(data.times)) * 0.95)
	data.nineTime = data.times[index]

}

func TailNginx(nginx Nginx, graphite Graphite, report Report, hostname string) {

	var uriReportRegex *regexp.Regexp
	if len(report.UriRegex) > 0 {
		uriReportRegex = regexp.MustCompile(report.UriRegex)
		if uriReportRegex == nil {
			log.Fatal("Unable to compile regex, please ensure this is a valid regular expression")
		}
	}

	var mutex = &sync.Mutex{}
	var data = new(nginxData)
	data.reset()

	seek := &tail.SeekInfo{0, 2}

	if conn, addr, err := connectToGraphite(graphite.Server); err == nil {
		duration := time.Duration(graphite.Interval) * time.Second
		label := hostname
		if report.Label != "" {
			label = fmt.Sprintf("%v.%v", hostname, report.Label)
		}
		go sendAtInterval(duration, label, conn, addr, data, mutex)

	} else {
		log.Panicf("Error connecting to Graphite server: %v", err)
	}

	for {
		if t, err := tail.TailFile(nginx.Logfile, tail.Config{Follow: true, Location: seek}); err == nil {
			for rawLine := range t.Lines {
				line := parseNginxLogLine(rawLine.Text)

				// Filter out by report parameters
				if report.Host != "" && line.host != report.Host {
					continue
				} else if len(report.Statuses) > 0 {
					found := false
					for _, s := range report.Statuses {
						if line.status == s {
							found = true
						}
					}
					if !found {
						continue
					}
				} else if len(report.Methods) > 0 {
					found := false
					for _, m := range report.Methods {
						if line.method == m {
							found = true
						}
					}
					if !found {
						continue
					}
				} else if report.Upstream != "" && line.upstream != report.Upstream {
					continue
				} else if report.Upstream != "" && line.upstream != report.Upstream {
					continue
				} else if report.UriRegex != "" {
					if !uriReportRegex.MatchString(line.uri) {
						continue
					}
				}

				mutex.Lock()
				data.registerHit(line)
				mutex.Unlock()
			}

		} else {
			log.Fatal("Error tailing file: ", err)
		}
	}
}

func parseNginxLogLine(text string) (line nginxLogLine) {

	// Request Time
	if r := timeRegex.FindAllStringSubmatch(text, -1); len(r) >= 1 {
		if x, err := strconv.ParseFloat(r[0][1], 64); err == nil {
			line.time = x
		}
	}
	// HTTP Code
	if r := statusRegex.FindAllStringSubmatch(text, -1); len(r) >= 1 {
		if x, err := strconv.ParseInt(r[0][1], 10, 64); err == nil {
			line.status = x
		}
	}
	// SSL
	if r := sslRegex.FindAllStringSubmatch(text, -1); len(r) >= 1 {
		line.ssl = r[0][1]
	}
	// Method
	if r := methodRegex.FindAllStringSubmatch(text, -1); len(r) >= 1 {
		line.method = r[0][1]
	}
	// Host
	if r := hostRegex.FindAllStringSubmatch(text, -1); len(r) >= 1 {
		line.host = r[0][1]
	}
	// Upstream
	if r := upstreamRegex.FindAllStringSubmatch(text, -1); len(r) >= 1 {
		line.upstream = r[0][1]
	}
	// URI
	if r := uriRegex.FindAllStringSubmatch(text, -1); len(r) >= 1 {
		line.uri = r[0][1]
	}

	return line
}

func sendAtInterval(interval time.Duration, label string, conn *net.UDPConn, addr *net.UDPAddr, data *nginxData, mutex *sync.Mutex) {
	for {
		time.Sleep(interval)
		mutex.Lock()

		rps := float64(data.hits) / interval.Seconds()

		writeData(fmt.Sprintf("%v.rps", label), rps, conn, addr)
		writeData(fmt.Sprintf("%v.normal", label), float64(data.percentNormal), conn, addr)
		writeData(fmt.Sprintf("%v.warn", label), float64(data.percentWarn), conn, addr)
		writeData(fmt.Sprintf("%v.error", label), float64(data.percentError), conn, addr)
		writeData(fmt.Sprintf("%v.min", label), data.minTime, conn, addr)
		writeData(fmt.Sprintf("%v.max", label), data.maxTime, conn, addr)
		writeData(fmt.Sprintf("%v.avg", label), data.avgTime, conn, addr)
		writeData(fmt.Sprintf("%v.nine", label), data.nineTime, conn, addr)

		data.reset()
		mutex.Unlock()
	}
}
