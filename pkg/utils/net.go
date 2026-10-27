package utils

import (
	"net"

	"github.com/stashapp/stash/pkg/sliceutil"
)

func shouldAdvertiseIP(ip net.IP) bool {
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	return !ip4.IsLoopback() && !ip4.IsUnspecified() && !ip4.IsLinkLocalUnicast()
}

// preferredOutboundIPv4 returns the address the OS would use to reach the internet; dialing UDP sends no packet.
func preferredOutboundIPv4() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return nil
	}

	return addr.IP
}

// LocalIPv4s returns the IPv4 addresses other devices on the LAN can reach this host on, preferred outbound address first.
func LocalIPv4s() []string {
	ips := make([]string, 0)

	add := func(ip net.IP) {
		if !shouldAdvertiseIP(ip) {
			return
		}

		ips = sliceutil.AppendUnique(ips, ip.To4().String())
	}

	add(preferredOutboundIPv4())

	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				add(v.IP)
			case *net.IPAddr:
				add(v.IP)
			}
		}
	}

	return ips
}
