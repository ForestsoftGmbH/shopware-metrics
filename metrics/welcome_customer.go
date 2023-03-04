package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type WelcomeCustomerCount struct {
	Counter *prometheus.GaugeVec
}

func NewWelcomeCounter() WelcomeCustomerCount {
	orderCountMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shopware_welcome_count",
			Help: "Number of existing customers with a setted password",
		},
		[]string{"sales_channel"},
	)
	orderCount := WelcomeCustomerCount{
		Counter: orderCountMetrics,
	}
	return orderCount
}
func (o WelcomeCustomerCount) GetGauge() *prometheus.GaugeVec {
	return o.Counter
}
func (o WelcomeCustomerCount) Grab(db *sql.DB) (*prometheus.GaugeVec, error) {
	var orderCountMetrics = o.Counter

	var orderCount int
	//iterate sales channels

	sql := "SELECT COUNT(*) FROM `boa_data` WHERE welcome_hash IS NULL"
	err2 := db.QueryRow(sql).Scan(&orderCount)
	if err2 != nil {
		log.Println("Error", err2, sql)
	} else {
		orderCountMetrics.WithLabelValues("Boa Metal Solutions GmbH").Set(float64(orderCount))
	}
	return orderCountMetrics, nil
}
