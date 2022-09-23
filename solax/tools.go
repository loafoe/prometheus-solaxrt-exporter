package solax

import (
	"fmt"
	"net"
	"net/netip"
	"net/url"
)

func LocallyReachable(apiAddr string) (bool, error) {
	addr, err := url.Parse(apiAddr)
	if err != nil {
		return false, err
	}
	ip, err := netip.ParseAddr(addr.Hostname())
	if err != nil {
		return false, fmt.Errorf("expecting IP address, not host")
	}

	list, err := net.Interfaces()
	if err != nil {
		return false, err
	}

	for _, iface := range list {
		addrs, err := iface.Addrs()
		if err != nil {
			return false, err
		}
		for _, addr := range addrs {
			network, err := netip.ParsePrefix(addr.String())
			if err != nil {
				continue
			}
			if network.Contains(ip) {
				return true, nil
			}
		}
	}
	return false, nil
}
