package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"shopware-metrics/metrics"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/utils/env"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", *user, *password, *host, *dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	orderCountMetrics, err := metrics.NewOrderCount(*db)
	if err != nil {
		log.Println(err)
	}

	reg.MustRegister(orderCountMetrics)

	log.Printf("Starting Server at %s/metrics", *addr)
	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(*addr, nil))

}
