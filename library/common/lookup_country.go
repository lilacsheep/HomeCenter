package common

import (
	"errors"
	"net"
	"regexp"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/oschwald/geoip2-golang"
)

const (
	apiUrl string = "http://api.dukeshi.cn:29090/api/location"
)

var (
	GeoFileNotFoundErr = errors.New("")
	GeoDB              *geoip2.Reader
)

func init() {
	var err error
	GeoDB, err = geoip2.Open("data/GeoLite2-Country.mmdb")
	if err != nil {
		panic(err)
	}
}

func LookupByApi(addr string) (string, error) {
	client := ghttp.NewClient()
	resp, err := client.Get(apiUrl, g.Map{"ip": addr})
	if err != nil {
		return "", err
	}
	defer resp.Close()
	data := resp.ReadAllString()
	j := gjson.New(data)
	return j.GetString("country"), nil
}

func LookupCountry(addr string) (string, error) {
	record, err := GeoDB.Country(net.ParseIP(addr))
	if err != nil {
		glog.Warning("get ip from geo data error: " + err.Error())
		return LookupByApi(addr)
	}
	return record.Country.IsoCode, err
}

func CheckIp(ip string) bool {
	addr := strings.Trim(ip, " ")
	regStr := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	if match, _ := regexp.MatchString(regStr, addr); match {
		return true
	}
	return false
}
