package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/loafoe/prometheus-solaxrt-exporter/solax"
	"github.com/loafoe/prometheus-solaxrt-exporter/solax/inverter"
	"github.com/loafoe/prometheus-solaxrt-exporter/solax/inverter/fields"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var listenAddr string
var apiAddr string
var debug bool

var (
	metricNamePrefix = "solaxrt_"
	registry         = prometheus.NewRegistry()
)

var (
	yieldTodayMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "yield_today",
		Help: "The yield for today (KWh)",
	}, []string{
		"inverter_sn",
	})

	yieldTotalMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "yield_total",
		Help: "The total yield of the system (KWh)",
	}, []string{
		"inverter_sn",
	})

	acPowerMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "ac_power",
		Help: "Current power generation (Wh)",
	}, []string{
		"inverter_sn",
	})
	upMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricNamePrefix + "up",
		Help: "The inverter power on status",
	}, []string{
		"sn",
	})
)

func init() {
	registry.MustRegister(yieldTotalMetrics)
	registry.MustRegister(yieldTodayMetric)
	registry.MustRegister(acPowerMetric)
	registry.MustRegister(upMetric)
}

func main() {
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
	flag.StringVar(&listenAddr, "listen", "0.0.0.0:8886", "Listen address for HTTP metrics")
	flag.StringVar(&apiAddr, "address", "http://5.8.8.8", "The address of the Realtime Inverter interface")
	flag.Parse()

	go func() {
		sleep := false
		for {
			if sleep {
				time.Sleep(time.Second * 5)
			}
			sleep = true
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

			fmt.Printf("calling Realtime API...\n")
			resp, err := solax.GetRealtimeInfo[inverter.X1BoostAirMini](ctx,
				solax.WithURL(apiAddr),
				solax.WithDebug(debug))
			cancel()
			if err != nil {
				fmt.Printf("error: %v\n", err)
				upMetric.WithLabelValues("").Set(0)
				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Printf("not sleeping\n")
					sleep = false
				}
				continue
			}
			yieldTodayMetric.WithLabelValues(resp.SN).Set(resp.Field(fields.Todays_Energy))
			yieldTotalMetrics.WithLabelValues(resp.SN).Set(resp.Field(fields.Total_Energy))
			acPowerMetric.WithLabelValues(resp.SN).Set(resp.Field(fields.AC_Power))
			upMetric.WithLabelValues("").Set(1.0)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	_ = http.ListenAndServe(listenAddr, nil)
}
