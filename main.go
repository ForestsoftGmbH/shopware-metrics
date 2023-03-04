package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/utils/env"
	"log"
	"net/http"
	"os"
	"shopware-metrics/database"
	"shopware-metrics/metrics"
	"time"
)

var addr = flag.String("listen-address", "0.0.0.0:8090", "The address to listen on for HTTP requests.")

func main() {
	user := flag.String("user", env.GetString("DB_USER", "root"), "user")
	password := flag.String("password", env.GetString("DB_PASSWORD", "root"), "password")
	host := flag.String("host", env.GetString("DB_HOST", "localhost"), "host")
	dbname := flag.String("dbname", env.GetString("DB_NAME", "shopware"), "dbname")

	flag.Parse()

	// Create non-global registry.
	reg := prometheus.NewRegistry()
	config := database.DbConfig{
		User:     *user,
		Password: *password,
		Host:     *host,
		Dbname:   *dbname,
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	orderCounter := metrics.NewOrderCount(config)
	orderCountMetrics, err := orderCounter.Grab()
	if err != nil {
		fmt.Println(err)
	} else {
		reg.MustRegister(orderCountMetrics)
	}

	log.Printf("Starting Server at %s/metrics", *addr)
	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	go func() {
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()

	defer func() {
		cancel()
	}()

	//make a new slice of metrics
	metrics := []metrics.ShopwareMetrics{
		orderCounter,
	}
	if err := run(ctx, metrics); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
func run(ctx context.Context, metrics []metrics.ShopwareMetrics) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(time.Duration(10) * time.Second):
			for _, metric := range metrics {
				log.Println("Grabbing metrics for %v", metric)
				_, err := metric.Grab()
				if err != nil {
					log.Println("Error: ", err)
				}
			}
		}
	}
	return nil
}
