package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
)

type OrderCount struct {
	Counter *prometheus.GaugeVec
}

type ShopwareMetrics interface {
	Grab(db *sql.DB) (*prometheus.GaugeVec, error)
}
