package ddns

import (
	"github.com/gogf/gf/net/ghttp"
	"net"
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

func (self Address) RecordType() RecordType {
	if self.IsIpv4() {
		return AType
	}
	if self.IsIpv6() {
		return AAAAType
	}
	return ""
}

func (self Address) IsIpv4() bool {
	return strings.Count(string(self), ":") < 2
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

func GetPublicAddress() (Address, error) {
	url := "https://www.dukeshi.com/api/ip"
	cli := ghttp.NewClient()
	resp, err := cli.Get(url)
	if err != nil {
		return "", err
	} else {
		return Address(resp.ReadAllString()), nil
	}
}

func GetPublicV6AddressByInterface(index int) ([]Address, error) {
	var (
		Addresses []Address
		err       error
		inter     *net.Interface
	)
	inter, err = net.InterfaceByIndex(index)
	if err != nil {
		return nil, err
	}
	if inter.Flags&net.FlagUp == 1 {
		if inter.Flags&net.FlagLoopback == 0 {
			addrs, _ := inter.Addrs()
			for _, a := range addrs {
				address := Address(a.String())
				if address.Verify() {
					if address.IsPublic() {
						Addresses = append(Addresses, address)
					}
				}
			}
		}
	}
	return Addresses, err
}
