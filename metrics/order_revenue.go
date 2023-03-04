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
		[]string{"sales_channel"},
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
	var orderCount int

	for _, salesChannel := range salesChannels {
		sql := "SELECT COUNT(*) FROM `order` WHERE sales_channel_id = ?"
		err2 := db.QueryRow(sql, salesChannel.Id).Scan(&orderCount)
		if err2 != nil {
			log.Println("Error", err2, sql)
			log.Println("Sales Channel:", salesChannel.Id)
			log.Println("Sales Channel Name:", salesChannel)
		} else {
			orderCountMetrics.WithLabelValues(salesChannel.Name).Set(float64(orderCount))
		}
	}
	return orderCountMetrics, nil
}
