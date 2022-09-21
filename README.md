# prometheus-solaxrt-exporter

Prometheus exporter for Realtime Solax Inverter readouts

## Install

Using Go 1.19 or newer

```shell
go install github.com/loafoe/prometheus-solaxrt-exporter@latest
```

## Usage

### Run exporter

```shell
prometheus-solaxcloud-exporter -listen 0.0.0.0:8886
```

### Ship to prometheus

You can use something like Grafana-agent to ship data to a remote write endpoint. Example:

```yml
metrics:
  configs:
    - name: default
      scrape_configs:
        - job_name: 'solaxcloud_exporter'
          scrape_interval: 2m
          static_configs:
            - targets: ['localhost:8886']
      remote_write:
        - url: https://prometheus.example.com/api/v1/write
          basic_auth:
            username: scraper
            password: S0m3pAssW0rdH3Re
```

## License

License is MIT
