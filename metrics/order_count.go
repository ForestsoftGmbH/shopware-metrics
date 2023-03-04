package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"shopware-metrics/database"
)

type OrderCount struct {
	Counter  *prometheus.CounterVec
	dbconfig database.DbConfig
}

func NewOrderCount(dp database.DbConfig) OrderCount {
	orderCountMetrics := prometheus.NewCounterVec(
		prometheus.CounterOpts{
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

func (o OrderCount) Grab() (*prometheus.CounterVec, error) {
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
		err2 := db.QueryRow("SELECT COUNT(*) FROM `order` WHERE sales_channel_id = ?", channelId).Scan(&orderCount)
		if err2 != nil {
			log.Fatal(err2)
		}

		orderCountMetrics.WithLabelValues(salesChannel).Add(float64(orderCount))
	}
	return orderCountMetrics, err
}
