package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"shopware-metrics/shopware"
)

type OrderCount struct {
	Counter *prometheus.GaugeVec
}

func NewOrderCount() OrderCount {
	orderCountMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shopware_order_count",
			Help: "Number of orders",
		},
		[]string{"sales_channel", "duration"},
	)
	orderCount := OrderCount{
		Counter: orderCountMetrics,
	}
	return orderCount
}
func (o OrderCount) GetGauge() *prometheus.GaugeVec {
	return o.Counter
}
func (o OrderCount) Grab(db *sql.DB) (*prometheus.GaugeVec, error) {
	salesChannels := shopware.GetSalesChannels(db)
	for _, salesChannel := range salesChannels {
		o.grabDatabase(db, salesChannel, "daily")
		o.grabDatabase(db, salesChannel, "hourly")
		o.grabDatabase(db, salesChannel, "total")
	}
	return o.Counter, nil
}

func (o OrderCount) grabDatabase(db *sql.DB, salesChannel shopware.SalesChannel, metric string) {
	var orderCount int
	var interval string

	switch metric {
	case "daily":
		interval = "AND created_at > DATE_SUB(NOW(), INTERVAL 1 DAY)"
	case "hourly":
		interval = "AND created_at > DATE_SUB(NOW(), INTERVAL 1 HOUR)"
	case "total":
		interval = ""
	}

	sql := "SELECT COUNT(*) FROM `order` WHERE sales_channel_id = ? " + interval
	err2 := db.QueryRow(sql, salesChannel.Id).Scan(&orderCount)
	if err2 != nil {
		log.Println("Error", err2, sql)
		log.Println("Sales Channel Name:", salesChannel.Name)
	} else {
		o.Counter.WithLabelValues(salesChannel.Name, metric).Set(float64(orderCount))
	}
}
