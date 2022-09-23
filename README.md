# prometheus-solaxrt-exporter

Prometheus exporter for real-time Solax Inverter data readouts

## Usage

This exporter is only tested with an X1-Boost-Mini Inverter with a Pocket Wifi dongle. The host should have a direct Wifi connection to the Pocket Wifi network. In practice, this means you'll probably want a dedicated compute module (Raspberry Pi) connected to this network. The default `http://5.8.8.8` hardcoded IP address is unfortunately publicly routable (pointing to a host in Russia of all places!). 

> The exporter contains code which ensures direct connectivity to the Pocket Wifi before attempting to query the real-time API endpoint.

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
          scrape_interval: 2s
          static_configs:
            - targets: ['localhost:8886']
      remote_write:
        - url: https://prometheus.example.com/api/v1/write
          basic_auth:
            username: scraper
            password: S0m3pAssW0rdH3Re
```

## Acknowledgement

API field mappings discovered from project https://github.com/squishykid/solax -- kudos!

## License

License is MIT
