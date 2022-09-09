package qnet

import "net"

func ToIpHostPort(hostPort string) (string, error) {
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", err
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(ips[0].String(), port), nil
}
