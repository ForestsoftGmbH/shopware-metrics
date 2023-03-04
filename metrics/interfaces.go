package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
)

type ShopwareMetrics interface {
	GetGauge() *prometheus.GaugeVec
	Grab(db *sql.DB) (*prometheus.GaugeVec, error)
}
