package solax_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loafoe/prometheus-solaxrt-exporter/solax"
	"github.com/loafoe/prometheus-solaxrt-exporter/solax/inverter"
	"github.com/stretchr/testify/assert"
)

var (
	muxSolaxCloud    *http.ServeMux
	serverSolaxCloud *httptest.Server
)

func setup(_ *testing.T) func() {
	muxSolaxCloud = http.NewServeMux()
	serverSolaxCloud = httptest.NewServer(muxSolaxCloud)

	return func() {
		serverSolaxCloud.Close()
	}
}

func getURL(host string) string {
	return host + "/"
}

func TestGetRealtimeInfo(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxSolaxCloud.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, http.MethodPost, r.Method) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{
  "type": "X1-Boost-Air-Mini",
  "SN": "S123456789",
  "ver": "2.32.6",
  "Data": [
    0.3,
    0,
    142.9,
    0,
    0.5,
    238.5,
    61,
    33,
    6.5,
    5665.1,
    0,
    71,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    50,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    2
  ],
  "Information": [
    1.5,
    4,
    "X1-Boost-Air-Mini",
    "XM123456789101112",
    1,
    3.21,
    1.08,
    1.1,
    0
  ]
}`)
	})

	ctx := context.Background()

	resp, err := solax.GetRealtimeInfo[inverter.X1BoostAirMini](ctx, solax.WithURL(getURL(serverSolaxCloud.URL)))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Equal(t, "S123456789", resp.SN)
}
