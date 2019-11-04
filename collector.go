package main

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	connectionEstablished = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "drbd_connection_established",
			Help: "a",
		},
		[]string{"resource", "connection", "id"},
	)
	diskUpToDate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "drbd_local_disk_up_to_date",
			Help: "a",
		},
		[]string{"resource", "connection", "id"},
	)
	connectionUnhealthy = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "drbd_connection_unhealthy",
			Help: "a",
		},
		[]string{"resource", "connection", "id"},
	)
	gauges = map[string]*prometheus.GaugeVec{
		"ns": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_network_sent_bytes",
				Help: "Volume of net data sent to the partner via the network connection.",
			},
			[]string{"resource", "connection", "id"}),
		"nr": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_network_received_bytes",
				Help: "Volume of net data received by the partner via the network connection.",
			},
			[]string{"resource", "connection", "id"}),
		"dw": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_disk_written_bytes",
				Help: "Net data written on local hard disk.",
			},
			[]string{"resource", "connection", "id"}),
		"dr": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_disk_read_bytes",
				Help: "Net data read from local hard disk.",
			},
			[]string{"resource", "connection", "id"}),
		"al": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_activitylog_writes",
				Help: "Number of updates of the activity log area of the meta data.",
			},
			[]string{"resource", "connection", "id"}),
		"bm": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_bitmap_writes",
				Help: "Number of updates of the bitmap area of the meta data.",
			},
			[]string{"resource", "connection", "id"}),

		"lo": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_local_pending",
				Help: "Number of open requests to the local I/O sub-system.",
			},
			[]string{"resource", "connection", "id"}),
		"pe": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_remote_pending",
				Help: "Number of requests sent to the partner, but that have not yet been answered by the latter.",
			},
			[]string{"resource", "connection", "id"}),
		"ua": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_remote_unacknowledged",
				Help: "Number of requests received by the partner via the network connection, but that have not yet been answered.",
			},
			[]string{"resource", "connection", "id"}),
		"ap": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_application_pending",
				Help: "Number of block I/O requests forwarded to DRBD, but not yet answered by DRBD.",
			},
			[]string{"resource", "connection", "id"}),
		"ep": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_epochs",
				Help: "Number of Epochs currently on the fly.",
			},
			[]string{"resource", "connection", "id"}),
		"oos": prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "drbd_out_of_sync_bytes",
				Help: "Amount of data known to be out of sync.",
			},
			[]string{"resource", "connection", "id"}),
	}
)

//Collector implements the prometheus.Collector interface.
type Collector struct{}

//Describe implements the prometheus.Collector interface.
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	connectionEstablished.Describe(ch)
	diskUpToDate.Describe(ch)
	connectionUnhealthy.Describe(ch)
	for _, v := range gauges {
		v.Describe(ch)
	}
}

//Collect implements the prometheus.Collector interface.
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	err := c.measure()
	//only report data when measurement was successful
	if err == nil {
		connectionEstablished.Collect(ch)
		diskUpToDate.Collect(ch)
		connectionUnhealthy.Collect(ch)

		for _, v := range gauges {
			v.Collect(ch)
		}
	} else {
		log.Println("ERROR:", err)
		return
	}
}

var errNoStatus = fmt.Errorf("No DRBD volumes detected")

func (c Collector) measure() error {
	connections := getAllDRDBstatues()
	if len(connections) == 0 {
		return errNoStatus
	}

	for _, connection := range connections {

		isEstablished := (connection.ConnectionStatus == "Established" || connection.ConnectionStatus == "SyncSource" || connection.ConnectionStatus == "SyncTarget")
		if isEstablished {
			connectionEstablished.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(1)
			connectionUnhealthy.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(0)
		} else {
			connectionEstablished.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(0)
			connectionUnhealthy.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(1)
		}

		if connection.MyDiskStatus != "UpToDate" {
			diskUpToDate.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(0)
		} else {
			diskUpToDate.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(1)
		}

		for _, kv := range connection.KVs {
			if gauges[kv.Name] != nil {
				gauges[kv.Name].WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(kv.Value)
			}
		}

	}
	return nil
}
