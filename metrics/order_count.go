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
		[]string{"sales_channel"},
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
	var orderCountMetrics = o.Counter
	salesChannels := shopware.GetSalesChannels(db)
	var orderCount int

	for _, salesChannel := range salesChannels {
		sql := "SELECT COUNT(*) FROM `order` WHERE sales_channel_id = ?"
		err2 := db.QueryRow(sql, salesChannel.Id).Scan(&orderCount)
		if err2 != nil {
			log.Println("Error", err2, sql)
			log.Println("Sales Channel Name:", salesChannel.Name)
		} else {
			orderCountMetrics.WithLabelValues(salesChannel.Name).Set(float64(orderCount))
		}
	}
	return orderCountMetrics, nil
}
