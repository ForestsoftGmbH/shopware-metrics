package shopware

import (
	"database/sql"
	"log"
)

type SalesChannel struct {
	Id   string
	Name string
}

func GetSalesChannels(db *sql.DB) []SalesChannel {
	res, err := db.Query("SELECT sales_channel.id, sales_channel_translation.name FROM sales_channel INNER JOIN sales_channel_translation ON sales_channel.id = sales_channel_translation.sales_channel_id")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()
	var salesChannels []SalesChannel
	for res.Next() {
		var salesChannel SalesChannel
		err := res.Scan(&salesChannel.Id, &salesChannel.Name)
		if err != nil {
			log.Fatal(err)
		}
		salesChannels = append(salesChannels, salesChannel)
	}
	return salesChannels
}
