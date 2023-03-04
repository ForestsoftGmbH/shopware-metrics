package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"shopware-metrics/shopware"
)

type CustomerGauge struct {
	Counter *prometheus.GaugeVec
}

func NewCustomerCounter() CustomerGauge {
	orderCountMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shopware_customer_count",
			Help: "Number of customers",
		},
		[]string{"sales_channel"},
	)
	orderCount := CustomerGauge{
		Counter: orderCountMetrics,
	}
	return orderCount
}
func (o CustomerGauge) GetGauge() *prometheus.GaugeVec {
	return o.Counter
}
func (o CustomerGauge) Grab(db *sql.DB) (*prometheus.GaugeVec, error) {
	var orderCountMetrics = o.Counter
	salesChannels := shopware.GetSalesChannels(db)
	var orderCount int
	//iterate sales channels
	for _, salesChannel := range salesChannels {
		sql := "SELECT COUNT(*) FROM `customer` WHERE sales_channel_id = ?"
		err2 := db.QueryRow(sql, salesChannel.Id).Scan(&orderCount)
		if err2 != nil {
			log.Println("Error", err2, sql)
			log.Println("Sales Channel:", salesChannel.Id)
			log.Println("Sales Channel Name:", salesChannel.Name)
		} else {
			orderCountMetrics.WithLabelValues(salesChannel.Name).Set(float64(orderCount))
		}
	}
	return orderCountMetrics, nil
}
