package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"shopware-metrics/shopware"
)

type RevenueGauge struct {
	Counter *prometheus.GaugeVec
}

func NewOrderRevenue() RevenueGauge {
	orderCountMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shopware_order_revenue_net",
			Help: "Number of orders",
		},
		[]string{"sales_channel", "metric"},
	)
	orderCount := RevenueGauge{
		Counter: orderCountMetrics,
	}
	return orderCount
}
func (o RevenueGauge) GetGauge() *prometheus.GaugeVec {
	return o.Counter
}
func (o RevenueGauge) Grab(db *sql.DB) (*prometheus.GaugeVec, error) {

	var orderCountMetrics = o.Counter
	salesChannels := shopware.GetSalesChannels(db)
	for _, salesChannel := range salesChannels {
		o.grabDatabase(db, salesChannel, "daily")
		o.grabDatabase(db, salesChannel, "hourly")
		o.grabDatabase(db, salesChannel, "total")
	}
	return orderCountMetrics, nil
}

func (o RevenueGauge) grabDatabase(db *sql.DB, salesChannel shopware.SalesChannel, metric string) {
	var orderCount float64
	var interval string

	switch metric {
	case "daily":
		interval = "AND created_at > DATE_SUB(NOW(), INTERVAL 1 DAY)"
	case "hourly":
		interval = "AND created_at > DATE_SUB(NOW(), INTERVAL 1 HOUR)"
	case "total":
		interval = ""
	}
	sql := "SELECT IFNULL(SUM(amount_net),0) FROM `order` WHERE sales_channel_id = ? " + interval
	err2 := db.QueryRow(sql, salesChannel.Id).Scan(&orderCount)
	if err2 != nil {
		log.Println("Error", err2, sql)
		log.Println("Sales Channel Name:", salesChannel.Name)
	} else {
		o.Counter.WithLabelValues(salesChannel.Name, metric).Set(orderCount)
	}
}
