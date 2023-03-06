package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"shopware-metrics/shopware"
)

type CartGauge struct {
	Counter *prometheus.GaugeVec
}

func NewCartCounter() CartGauge {
	orderCountMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shopware_cart_count",
			Help: "Sum of open carts",
		},
		[]string{"sales_channel", "metric", "timespan"},
	)
	orderCount := CartGauge{
		Counter: orderCountMetrics,
	}
	return orderCount
}
func (o CartGauge) GetGauge() *prometheus.GaugeVec {
	return o.Counter
}
func (o CartGauge) Grab(db *sql.DB) (*prometheus.GaugeVec, error) {
	var orderCountMetrics = o.Counter
	salesChannels := shopware.GetSalesChannels(db)
	//iterate sales channels
	for _, salesChannel := range salesChannels {

		o.grabDatabase(db, salesChannel, "daily", "revenue")
		o.grabDatabase(db, salesChannel, "hourly", "revenue")
		o.grabDatabase(db, salesChannel, "total", "revenue")

		o.grabDatabase(db, salesChannel, "daily", "line_item_count")
		o.grabDatabase(db, salesChannel, "hourly", "line_item_count")
		o.grabDatabase(db, salesChannel, "total", "line_item_count")

	}
	return orderCountMetrics, nil
}

func (o CartGauge) grabDatabase(db *sql.DB, salesChannel shopware.SalesChannel, metric string, label string) {
	var orderCount float64
	var interval string
	var value string

	switch metric {
	case "daily":
		interval = "AND created_at > DATE_SUB(NOW(), INTERVAL 1 DAY)"
	case "hourly":
		interval = "AND created_at > DATE_SUB(NOW(), INTERVAL 1 HOUR)"
	case "total":
		interval = ""
	}

	switch label {
	case "revenue":
		value = "IFNULL(SUM(price), 0)"
	case "line_item_count":
		value = "IFNULL(SUM(line_item_count), 0)"
	}

	sql := "SELECT " + value + " FROM `cart` WHERE sales_channel_id = ? " + interval
	err2 := db.QueryRow(sql, salesChannel.Id).Scan(&orderCount)
	if err2 != nil {
		log.Println("Error", err2, sql)
		log.Println("Sales Channel Name:", salesChannel.Name)
	} else {
		o.Counter.WithLabelValues(salesChannel.Name, label, metric).Set(orderCount)
	}
}
