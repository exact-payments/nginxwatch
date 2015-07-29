# Nginx Watch
A Go based daemon that tails an nginx logfile, generics metrics and sends them to a Graphite server.


# Installation

This app is designed to work and has been tested on Linux and OSX.

# Configuring nginx

In order to use this it expects to gather specific data from a specifically formatted nginx logfile.   You will need to add the following log definition for nginx:

    log_format stats_log "host:$host"
                         "\tupstream:$proxy_host"
                         "\tstatus:$upstream_status"
                         "\ttime:$upstream_response_time"
                         "\tssl:$ssl_protocol"
                         "\tmethod:$request_method"
                         "\turi:$uri";

Next you'll need to add to each of your server definition blocks an extra access_log statement for the above (modifying the path if needed):

    access_log  /var/log/nginx/stats.log stats_log;

Nginx will now log extra statistics to that log which nginxwatch can pick up and read.  You'll need to rotate this logfile.

# Running

Create a config.toml file with contents like the following:

		[nginx]
		logfile = '/var/log/nginx/status.log'

		[graphite]
		server = 'localhost:8999'
		interval = 10

Where interval is the seconds between submissions.  nginxwatch will send 8 points of graphite data every 10 seconds by default for the global statistics of the nginx logs.

# Primary Report

By default the report will spit out something like the following (you can see the output by using the debug flag -d):

  2015/07/29 17:25:05 proxy-1.telemetryapp.com.nine: 0.402
  2015/07/29 17:25:15 proxy-1.telemetryapp.com.rps: 5.5
  2015/07/29 17:25:15 proxy-1.telemetryapp.com.normal: 85
  2015/07/29 17:25:15 proxy-1.telemetryapp.com.warn: 1
  2015/07/29 17:25:15 proxy-1.telemetryapp.com.error: 0
  2015/07/29 17:25:15 proxy-1.telemetryapp.com.min: 0.013
  2015/07/29 17:25:15 proxy-1.telemetryapp.com.max: 2.591
  2015/07/29 17:25:15 proxy-1.telemetryapp.com.avg: 0.11638297872340424

nine, min, max, avg are the ninety-fifth percentile, minimum, maximum and average response times.  rps is the requests per second.  normal, warn and error are the hit counts of the normal (http 2xx), warning (http 4xx) and error (http 5xx) http codes. 

# Optional Additional Reports

You can configure nginxwatch to also watch for specific variables in the logfile and send an additional set of reports that just match this.  To do this create one or more [report] blocks as such:

    [[report]]
    label = 'api'
    upstream = 'api'
    host = 'api.telemetryapp.com'
    statuses = [200, 201]
    methods = ['GET', 'POST', 'PATCH', 'PUT', 'DELETE']
    uri_regex = ''

Where label is what Graphite will see.  Upstream is the optional nginx upstream to match.  Host is the optional virtual host to match.  Statuses are an optional array of statuses to match.  Methods are an optional array of methods to match.  URI_Regex is an optional regex to match against the URI.

# Building

nginxwatch is written in Go.  In order to build you'll need to have a working Go environment.  Typically we cross compile

To build for linux:

	GOOS=linux GOARCH=amd64 go build -o pkg/linux_amd64/nginxwatch

# Creating .deb

Usually you'll want to install from a .deb file rather than a binary.  To create a .deb you'll want to do the following:

Install the [FPM Ruby Gem](https://github.com/jordansissel/fpm)

	fpm -s dir -t deb -n 'nginxwatch' -v 1.0 -m support@telemetryapp.com --vendor support@telemetryapp.com --license MIT --url https://telemetryapp.com -e --prefix /bin -C pkg/linux_amd64 nginxwatch
