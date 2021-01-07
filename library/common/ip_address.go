package common

import (
	"net"
	"strings"
)


type IPAddress string

func (self IPAddress) ip() net.IP {
	return net.ParseIP(self.ip().String())
}

func (self IPAddress) String() string {
	c := strings.Index(string(self), "/")
	switch c {
	case 0:
		return string(self)
	default:
		t := strings.Split(string(self), "/")
		return t[0]
	}
}

func (self IPAddress) Verify() bool {
	address := self.ip()
	return address != nil
}

func (self IPAddress) IsIpv4() bool {
	ip := self.ip()
	if ip != nil {
		return ip.To4() != nil
	}
	return false
}

func (self IPAddress) IsIpv6() bool {
	return strings.Count(self.ip().String(), ":") >= 2
}

func (self IPAddress) IsPublic() bool {
	ip := self.ip()
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
    if ip4 := ip.To4(); ip4 != nil {
        switch true {
        case ip4[0] == 10:
            return false
        case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
            return false
        case ip4[0] == 192 && ip4[1] == 168:
            return false
        default:
            return true
        }
	}
	if self.IsIpv6() {
		index := []string{"fe80::"}
		for _, i := range index {
			if strings.HasPrefix(self.ip().String(), i) {
				return false
			}
		}
	}
	return true
}
