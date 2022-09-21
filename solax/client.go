package solax

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func GetRealtimeInfo[K any](ctx context.Context, opts ...OptionFunc) (*K, error) {
	client := resty.New()

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
