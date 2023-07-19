package util

import (
	"context"
	"net"
)

// LookupIPv4 resolves the given hostname to a IPv4.
// WARNING: It only returns the first IP even if the hostname resolves to several
func LookupIPv4(ctx context.Context, host string) (string, error) {
	var ip string
	if r := net.ParseIP(host); r != nil {
		ip = r.String()
	} else {
		// if it is not an IP, try to resolve it
		ips, err := net.DefaultResolver.LookupIP(ctx, "ip4", host)
		if err != nil {
			return "", err
		}

		// return only the firt IP returned
		ip = ips[0].String()
	}

	return ip, nil
}
