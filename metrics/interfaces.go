package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"shopware-metrics/database"
)

type OrderCount struct {
	Counter  *prometheus.GaugeVec
	dbconfig database.DbConfig
}

type ShopwareMetrics interface {
	Grab() (*prometheus.GaugeVec, error)
}
