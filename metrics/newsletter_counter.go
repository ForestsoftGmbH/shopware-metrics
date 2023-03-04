package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"shopware-metrics/shopware"
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
	var orderCountMetrics = o.Counter
	salesChannels := shopware.GetSalesChannels(db)
	var orderCount int
	var status string
	for _, salesChannel := range salesChannels {

		sql := "SELECT COUNT(*), `status` FROM `newsletter_recipient` WHERE sales_channel_id = ? GROUP BY `status`"
		err2 := db.QueryRow(sql, salesChannel.Id).Scan(&orderCount, &status)
		if err2 != nil && "sql: no rows in result set" != err2.Error() {
			if err2 != nil {
				log.Println("Error", err2, sql)
				log.Println("Sales Channel Name:", salesChannel)
			}
		} else {
			orderCountMetrics.WithLabelValues(salesChannel.Name, status).Set(float64(orderCount))
		}
	}
	return orderCountMetrics, nil
}
