package common

import (
	"net"
	"regexp"
	"strings"
)


type Address string

func (self Address) String() string {
	c := strings.Index(string(self), "/")
	switch c {
	case 0:
		return string(self)
	default:
		t := strings.Split(string(self), "/")
		return t[0]
	}
}

func (self Address) Verify() bool {
	if strings.Index(string(self), "/") > 0 {
		_, IPNet, _ := net.ParseCIDR(string(self))
		if c, b := IPNet.Mask.Size(); c == 64 && b == 128 {
			return true
		}
	} else {
		address := net.ParseIP(string(self))
		if address != nil {
			return true
		}
	}
	return false
}

func (self Address) IsIpv4() bool {
	addr := strings.Trim(string(self), " ")
	regStr := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	if match, _ := regexp.MatchString(regStr, addr); match {
		return true
	}
	return false
}

func (self Address) IsIpv6() bool {
	return strings.Count(string(self), ":") >= 2
}

func (self Address) IsPublic() bool {
	address := string(self)
	if self.IsIpv4() {
		index := []string{"192.168.", "172.16.", "10.", "127.0."}
		for _, i := range index {
			if strings.HasPrefix(address, i) {
				return false
			}
		}
	}
	if self.IsIpv6() {
		index := []string{"fe80::"}
		for _, i := range index {
			if strings.HasPrefix(address, i) {
				return false
			}
		}
	}
	return true
}
