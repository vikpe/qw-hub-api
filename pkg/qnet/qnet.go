package qnet

import "net"

func ToIpHostPort(address string) (string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", err
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(ips[0].String(), port), nil
}
