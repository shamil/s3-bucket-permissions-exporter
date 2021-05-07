package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const namespace = "s3_bucket_permissions"

var (
	listenAddress  = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Short('l').Default(":9199").String()
	metricsPath    = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	ignoredBuckets = kingpin.Flag("collector.ignored-buckets", "Regexp of buckets to ignore from collecting.").Default("^$").String()
)

func init() {
	prometheus.MustRegister(version.NewCollector(namespace))
}

func main() {
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("s3-bucket-permissions-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	prometheus.MustRegister(NewCollector())

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//nolint:errcheck
		w.Write([]byte(`<html>
			<head><title>Burrow Exporter</title></head>
			<body>
			<h1>Burrow Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Infoln("listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
