package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

func NewOrderCount(db sql.DB) (prometheus.CounterVec, error) {
	orderCountMetrics := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shopware_order_count",
			Help: "Number of orders",
		},
		[]string{"sales_channel"},
	)

	var salesChannel string
	var channelId string
	res, err := db.Query("SELECT sales_channel.id, sales_channel_translation.name FROM sales_channel INNER JOIN sales_channel_translation ON sales_channel.id = sales_channel_translation.sales_channel_id")
	defer res.Close()
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
	return *orderCountMetrics, err
}
