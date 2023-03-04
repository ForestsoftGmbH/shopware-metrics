package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type NewsletterSubscriberGauge struct {
	Counter *prometheus.GaugeVec
}

func NewNewsletterCounter() NewsletterSubscriberGauge {
	orderCountMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "shopware_newsletter_subscriber_count",
			Help: "Number of newsletter subscribers",
		},
		[]string{"sales_channel", "status"},
	)
	orderCount := NewsletterSubscriberGauge{
		Counter: orderCountMetrics,
	}
	return orderCount
}
func (o NewsletterSubscriberGauge) GetGauge() *prometheus.GaugeVec {
	return o.Counter
}
func (o NewsletterSubscriberGauge) Grab(db *sql.DB) (*prometheus.GaugeVec, error) {
	var salesChannel string
	var channelId string
	var orderCountMetrics = o.Counter
	res, err := db.Query("SELECT sales_channel.id, sales_channel_translation.name FROM sales_channel INNER JOIN sales_channel_translation ON sales_channel.id = sales_channel_translation.sales_channel_id")
	defer res.Close()
	var orderCount int
	var status string

	for res.Next() {
		err := res.Scan(&channelId, &salesChannel)
		if err != nil {
			log.Fatal(err)
		}
		sql := "SELECT COUNT(*), `status` FROM `newsletter_recipient` WHERE sales_channel_id = ? GROUP BY `status`"
		err2 := db.QueryRow(sql, channelId).Scan(&orderCount, &status)
		if "sql: no rows in result set" != err2.Error() {
			if err2 != nil {
				log.Println("Error", err2, sql)
				log.Println("Sales Channel Name:", salesChannel)
			} else {
				orderCountMetrics.WithLabelValues(salesChannel, status).Set(float64(orderCount))
			}
		}
	}
	return orderCountMetrics, err
}
