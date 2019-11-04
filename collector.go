package main

import (
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
)

//Collector implements the prometheus.Collector interface.
type Collector struct {
}

//Describe implements the prometheus.Collector interface.
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	connectionEstablished.Describe(ch)
	diskUpToDate.Describe(ch)
	connectionUnhealthy.Describe(ch)
}

//Collect implements the prometheus.Collector interface.
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	err := c.measure()
	//only report data when measurement was successful
	if err == nil {
		connectionEstablished.Collect(ch)
		diskUpToDate.Collect(ch)
		connectionUnhealthy.Collect(ch)
	} else {
		log.Println("ERROR:", err)
		return
	}
}

func (c Collector) measure() error {
	connections := getAllDRDBstatues()
	for _, connection := range connections {

		isEstablished := (connection.ConnectionStatus == "Established" || connection.ConnectionStatus == "SyncSource" || connection.ConnectionStatus == "SyncTarget")
		if isEstablished {
			connectionEstablished.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(1)
		} else {
			connectionEstablished.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(0)
			connectionUnhealthy.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(1)
		}

		if connection.MyDiskStatus != "UpToDate" {
			connectionEstablished.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(0)
		} else {
			connectionEstablished.WithLabelValues(connection.Resource, connection.RemoteHost, connection.ResourceID).Set(1)
		}

	}
	return nil
}
