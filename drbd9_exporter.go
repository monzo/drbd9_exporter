package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

var (
	listenAddress = flag.String("web.listen-address", ":9481", "Address on which to expose metrics and web interface")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	hostSysPath   = flag.String("host.sys-path", "/", "Path to host sysfs.")
)

func main() {
	flag.Parse()

	log.Print("Starting node_exporter", version.Info())
	log.Print("Build context", version.BuildContext())
	prometheus.MustRegister(Collector{})
	handler := promhttp.HandlerFor(prometheus.DefaultGatherer,
		promhttp.HandlerOpts{})

	http.Handle(*metricsPath, handler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>DRBD9 Exporter</title></head>
			<body>
			<h1>DRBD9 Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Print("Listening on", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatal(err)
	}
}
