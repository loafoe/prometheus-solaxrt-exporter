package solax

import (
	"github.com/go-resty/resty/v2"
)

type OptionFunc func(c *resty.Client, r *resty.Request) (*resty.Request, error)

func WithDefaultURL() OptionFunc {
	return func(c *resty.Client, r *resty.Request) (*resty.Request, error) {
		r.URL = "http://5.8.8.8"
		return r, nil
	}
}

func WithURL(url string) OptionFunc {
	return func(c *resty.Client, r *resty.Request) (*resty.Request, error) {
		r.URL = url
		return r, nil
	}
}

func WithDebug(debug bool) OptionFunc {
	return func(c *resty.Client, r *resty.Request) (*resty.Request, error) {
		c.SetDebug(debug)
		return r, nil
	}
}
