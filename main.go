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
var ethDevice string
var apiAddr string
var debug bool
var scrapeInterval time.Duration
var lastKnownSN string

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
	scrapeErrorsMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metricNamePrefix + "scrape_errors_total",
		Help: "Total number of scrape errors",
	}, []string{
		"reason",
	})
	scrapeSuccessMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: metricNamePrefix + "scrape_success_total",
		Help: "Total number of successful scrapes",
	})
	scrapeDurationMetric = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    metricNamePrefix + "scrape_duration_seconds",
		Help:    "Duration of scrape in seconds",
		Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 20},
	})
)

func init() {
	registry.MustRegister(yieldTotalMetrics)
	registry.MustRegister(yieldTodayMetric)
	registry.MustRegister(acPowerMetric)
	registry.MustRegister(upMetric)
	registry.MustRegister(scrapeErrorsMetric)
	registry.MustRegister(scrapeSuccessMetric)
	registry.MustRegister(scrapeDurationMetric)
}

func main() {
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
	flag.StringVar(&listenAddr, "listen", "0.0.0.0:8886", "Listen address for HTTP metrics")
	flag.StringVar(&apiAddr, "address", "http://5.8.8.8", "The address of the Realtime Inverter interface")
	flag.StringVar(&ethDevice, "device", "wlan0", "The ethernet device to check for Pocket wifi")
	flag.DurationVar(&scrapeInterval, "interval", 5*time.Second, "Scrape interval")
	flag.Parse()

	go func() {
		consecutiveErrors := 0
		maxBackoff := 30 * time.Second

		for {
			startTime := time.Now()

			ok, _ := solax.LocallyReachable(apiAddr)
			if !ok {
				fmt.Printf("address %s is not locally reachable, skipping refresh...\n", apiAddr)
				scrapeErrorsMetric.WithLabelValues("unreachable").Inc()
				consecutiveErrors++
				backoffSleep(consecutiveErrors, scrapeInterval, maxBackoff)
				continue
			}

			fmt.Printf("calling Realtime API at %s...\n", apiAddr)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			resp, err := solax.GetRealtimeInfo[inverter.X1BoostAirMini](ctx,
				solax.WithURL(apiAddr),
				solax.WithDebug(debug))
			cancel()

			scrapeDurationMetric.Observe(time.Since(startTime).Seconds())

			if err != nil {
				fmt.Printf("error: %v\n", err)
				upMetric.WithLabelValues("").Set(0)
				if lastKnownSN != "" {
					acPowerMetric.WithLabelValues(lastKnownSN).Set(0)
				}
				consecutiveErrors++

				if errors.Is(err, context.DeadlineExceeded) {
					scrapeErrorsMetric.WithLabelValues("timeout").Inc()
				} else {
					scrapeErrorsMetric.WithLabelValues("request").Inc()
				}

				backoffSleep(consecutiveErrors, scrapeInterval, maxBackoff)
				continue
			}

			consecutiveErrors = 0
			lastKnownSN = resp.SN
			scrapeSuccessMetric.Inc()
			yieldTodayMetric.WithLabelValues(resp.SN).Set(resp.Field(fields.Todays_Energy))
			yieldTotalMetrics.WithLabelValues(resp.SN).Set(resp.Field(fields.Total_Energy))
			acPowerMetric.WithLabelValues(resp.SN).Set(resp.Field(fields.AC_Power))
			upMetric.WithLabelValues("").Set(1.0)

			time.Sleep(scrapeInterval)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	_ = http.ListenAndServe(listenAddr, nil)
}

func backoffSleep(consecutiveErrors int, baseInterval, maxBackoff time.Duration) {
	backoff := baseInterval
	for i := 0; i < consecutiveErrors && backoff < maxBackoff; i++ {
		backoff = backoff * 2
	}
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	fmt.Printf("sleeping %v before next attempt (consecutive errors: %d)\n", backoff, consecutiveErrors)
	time.Sleep(backoff)
}
