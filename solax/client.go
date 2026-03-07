package solax

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	sharedClient *resty.Client
	clientOnce   sync.Once
)

func getClient() *resty.Client {
	clientOnce.Do(func() {
		transport := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          10,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			DisableKeepAlives:     false,
		}

		sharedClient = resty.New().
			SetTransport(transport).
			SetTimeout(15 * time.Second).
			SetRetryCount(3).
			SetRetryWaitTime(500 * time.Millisecond).
			SetRetryMaxWaitTime(2 * time.Second).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				if err != nil {
					return true
				}
				return r.StatusCode() >= 500
			})
	})
	return sharedClient
}

func GetRealtimeInfo[K any](ctx context.Context, opts ...OptionFunc) (*K, error) {
	client := getClient()

	request := client.R().SetQueryParams(map[string]string{
		"optType": "ReadRealTimeData",
	})
	request.Method = http.MethodPost
	request, _ = WithDefaultURL()(client, request)

	for _, o := range opts {
		r, err := o(client, request)
		if err != nil {
			return nil, err
		}
		request = r
	}
	resp, err := request.SetContext(ctx).Send()
	if err != nil {
		return nil, err
	}
	var jsonResponse K
	err = json.Unmarshal(resp.Body(), &jsonResponse)
	if err != nil {
		return nil, err
	}
	return &jsonResponse, nil
}
