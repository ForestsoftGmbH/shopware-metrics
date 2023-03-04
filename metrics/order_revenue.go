package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
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
	var salesChannel string
	var channelId string
	var orderCountMetrics = o.Counter
	res, err := db.Query("SELECT sales_channel.id, sales_channel_translation.name FROM sales_channel INNER JOIN sales_channel_translation ON sales_channel.id = sales_channel_translation.sales_channel_id")
	defer res.Close()
	var orderCount int

	for res.Next() {
		err := res.Scan(&channelId, &salesChannel)
		if err != nil {
			log.Fatal(err)
		}
		sql := "SELECT COUNT(*) FROM `order` WHERE sales_channel_id = ?"
		err2 := db.QueryRow(sql, channelId).Scan(&orderCount)
		if err2 != nil {
			log.Println("Error", err2, sql)
			log.Println("Sales Channel:", channelId)
			log.Println("Sales Channel Name:", salesChannel)
		} else {
			orderCountMetrics.WithLabelValues(salesChannel).Set(float64(orderCount))
		}
	}
	return orderCountMetrics, err
}
