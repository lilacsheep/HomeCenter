package common

import (
	"regexp"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func LookupCountry(addr string) (string, error) {
	client := ghttp.NewClient()
	resp, err := client.Get("https://www.dukeshi.com/api/location", g.Map{"ip": addr})
	if err != nil {
		return "", err
	}
	defer resp.Close()
	data := resp.ReadAllString()
	j := gjson.New(data)
	return j.GetString("country"), nil
}

func CheckIp(ip string) bool {
	addr := strings.Trim(ip, " ")
	regStr := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	if match, _ := regexp.MatchString(regStr, addr); match {
		return true
	}
	return false
}