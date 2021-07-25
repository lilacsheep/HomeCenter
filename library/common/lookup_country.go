package common

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	libgeo "github.com/nranchev/go-libGeoIP"
)

const (
	apiUrl     string = "http://api.dukeshi.cn:29090/api/location"
	geoDataUrl string = "https://raw.githubusercontent.com/Loyalsoldier/geoip/release/geoip.dat"
	dataFile   string = "data/geoip.dat"
)

var (
	gi *libgeo.GeoIP
	GeoFileNotFoundErr = errors.New("")
)

func LookupByGeo(addr string) (*libgeo.Location, error) {
	var err error
	if !gfile.Exists(dataFile) {
		return nil, GeoFileNotFoundErr
	}
	if gi == nil {
		gi, err = libgeo.Load(dataFile)
		if err != nil {
			return nil, err
		}
	}
	ip := gi.GetLocationByIP(addr)
	return ip, err
}

func InitGeoFile(proxyUrl ...string) error {
	cli := ghttp.NewClient()
	if len(proxyUrl) > 0 {
		cli.SetProxy(proxyUrl[0])
	}
	resp, err := cli.Get(geoDataUrl)
	if err != nil {
		return err
	}
	defer resp.Close()

	err = gfile.PutBytes("data/geoip.dat", resp.ReadAll())
	if err != nil {
		return err
	}
	gi, err = libgeo.Load(dataFile)
	if err != nil {
		return err
	}
	return nil
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
	return LookupByApi(addr)
	// location, err := LookupByGeo(addr)
	// if err != nil {
	// 	if err == GeoFileNotFoundErr {
	// 		go InitGeoFile()
	// 		return LookupByApi(addr)
	// 	} else {
	// 		return "", err
	// 	}
	// }
	// return location.CountryCode, nil
}

func CheckIp(ip string) bool {
	addr := strings.Trim(ip, " ")
	regStr := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	if match, _ := regexp.MatchString(regStr, addr); match {
		return true
	}
	return false
}
