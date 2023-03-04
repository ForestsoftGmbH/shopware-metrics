package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"shopware-metrics/database"
)

func NewOrderCount(dp database.DbConfig) OrderCount {
	orderCountMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shopware_order_count",
			Help: "Number of orders",
		},
		[]string{"sales_channel"},
	)

	orderCount := OrderCount{
		Counter:  orderCountMetrics,
		dbconfig: dp,
	}
	return orderCount
}

func (o OrderCount) Grab() (*prometheus.GaugeVec, error) {
	db, err := database.NewConnection(o.dbconfig)
	var salesChannel string
	var channelId string
	var orderCountMetrics = o.Counter
	res, err := db.Query("SELECT sales_channel.id, sales_channel_translation.name FROM sales_channel INNER JOIN sales_channel_translation ON sales_channel.id = sales_channel_translation.sales_channel_id")
	defer res.Close()
	defer db.Close()
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
