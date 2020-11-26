package ddns

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

const (
	DNSPodApiUrl string = "https://dnsapi.cn"
)

type RecordType string

const (
	AType     RecordType = "A"
	AAAAType  RecordType = "AAAA"
	CNameType RecordType = "CNAME"
	MXType    RecordType = "MX"
	TextType  RecordType = "TEXT"
	SrvType   RecordType = "SRV"
)

type Record struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Line      string `json:"line"`
	LineID    string `json:"line_id"`
	Type      string `json:"type"`
	TTL       string `json:"ttl"`
	Value     string `json:"value"`
	Weight    string `json:"weight"`
	MX        string `json:"mx"`
	Enable    string `json:"enable"`
	Status    string `json:"status"`
	UpdatedOn string `json:"updated_on"`
}

type Domain struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	PunyCode         string `json:"punycode"`
	Grade            string `json:"grade"`
	GradeTitle       string `json:"grade_title"`
	Status           string `json:"status"`     // enable disable
	ExtStatus        string `json:"ext_status"` // notexist dnserror "" 正常
	Records          string `json:"records"`    // 域名下记录总条数
	GroupID          string `json:"group_id"`   // 域名分组 ID
	IsMark           string `json:"is_mark"`    // 是否设置域名星标
	Remark           string `json:"remark"`
	IsVip            string `json:"is_vip"`            // yes or no
	SearchEnginePush string `json:"searchengine_push"` // yes or no
	UserID           string `json:"user_id"`
	CreatedOn        string `json:"created_on"`
	UpdatedOn        string `json:"updated_on"`
	TTL              string `json:"ttl"`
	CNameSpeedUp     string `json:"cname_speedup"` // disable or enable
	Owner            string `json:"owner"`
	DnsPodNS         string `json:"dnspod_ns"`
	IsBeian          string `json:"is_beian"`
}

type DNSPodResponse struct {
	Domains []Domain `json:"domains,omitempty"`
	Domain  Domain   `json:"domain,omitempty"`
	Status  struct {
		Code     int    `json:"code"`
		CreateAt string `json:"create_at"`
		Message  string `json:"message"`
	} `json:"status"`
	Records []Record `json:"records,omitempty"`
	Record  Record   `json:"record,omitempty"`
}

type DNSPod struct {
	ID           string // 用于鉴权的 API Token
	Token        string // 用于鉴权的 API Token
	Format       string // 返回的数据格式，可选，默认为xml，建议用json
	Lang         string // 返回的错误语言，可选，默认为en，建议用cn
	ErrorOnEmpty string // 没有数据时是否返回错误，可选，默认为yes，建议用no
}

func (self *DNSPod) Version() (interface{}, error) {
	cli := ghttp.NewClient()
	response, err := cli.Post(fmt.Sprintf("%s/Info.Version", DNSPodApiUrl), self.PublicParam())
	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		body := response.ReadAllString()
		resp := DNSPodResponse{}
		err := gjson.New(body).ToStruct(&resp)
		return resp, err
	}
}

func (self *DNSPod) PublicParam() g.Map {
	return g.Map{
		"login_token":    fmt.Sprintf("%s,%s", self.ID, self.Token),
		"format":         self.Format,
		"lang":           self.Lang,
		"error_on_empty": "yes",
	}
}

func (self *DNSPod) DomainList(keyword ...string) (interface{}, error) {
	cli := ghttp.NewClient()
	params := self.PublicParam()
	if len(keyword) != 0 {
		params["keyword"] = keyword[0]
	}
	response, err := cli.Post(fmt.Sprintf("%s/Domain.List", DNSPodApiUrl), params)
	if err != nil {
		return nil, err
	} else {
		defer response.Close()
		resp := DNSPodResponse{}
		err := gjson.New(response.ReadAllString()).ToStruct(&resp)
		return resp, err
	}
}

func (self *DNSPod) RecordList(domain string, subDomain string, keyword ...string) (interface{}, error) {
	cli := ghttp.NewClient()
	params := self.PublicParam()
	params["domain"] = domain
	if subDomain != "" {
		params["sub_domain"] = subDomain
	}
	if len(keyword) > 0 {
		params["keyword"] = keyword[0]
	}
	response := DNSPodResponse{}
	resp, err := cli.Post(fmt.Sprintf("%s/Record.List", DNSPodApiUrl), params)
	if err != nil {
		return nil, err
	} else {
		defer resp.Close()
		err = gjson.New(resp.ReadAllString()).ToStruct(&response)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
}

func (self *DNSPod) RecordCreate(domain, subDomain, ttl string, address Address) (string, error) {
	cli := ghttp.NewClient()
	params := self.PublicParam()
	params["domain"] = domain
	params["record_type"] = address.RecordType()
	if subDomain != "" {
		params["sub_domain"] = subDomain
	}
	params["value"] = address.String()
	params["record_line"] = "默认"
	if ttl != "" {
		params["ttl"] = ttl
	}
	resp, err := cli.Post(fmt.Sprintf("%s/Record.Create", DNSPodApiUrl), params)
	if err != nil {
		return "", err
	}
	defer resp.Close()
	response := DNSPodResponse{}
	err = gjson.New(resp.ReadAllString()).ToStruct(&response)
	if response.Status.Code != 1 {
		return "", errors.New(response.Status.Message)
	}
	return response.Record.ID, err
}

func (self *DNSPod) RecordRemove(domain string, recordID string) (interface{}, error) {
	cli := ghttp.NewClient()
	params := self.PublicParam()
	params["domain"] = domain
	params["record_id"] = recordID
	resp, err := cli.Post(fmt.Sprintf("%s/Record.Remove", DNSPodApiUrl), params)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	response := DNSPodResponse{}
	err = gjson.New(resp.ReadAllString()).ToStruct(&response)
	return response, err
}

func (self *DNSPod) RecordModify(domain, recordID, subDomain string, address Address) error {
	cli := ghttp.NewClient()
	params := self.PublicParam()
	params["domain"] = domain
	params["record_id"] = recordID
	params["value"] = address.String()
	params["record_line"] = "默认"
	if subDomain != "" {
		params["sub_domain"] = subDomain
	}
	params["record_type"] = address.RecordType()
	resp, err := cli.Post(fmt.Sprintf("%s/Record.Modify", DNSPodApiUrl), params)
	if err != nil {
		return err
	}
	defer resp.Close()
	response := DNSPodResponse{}
	err = gjson.New(resp.ReadAllString()).ToStruct(&response)
	if err != nil {
		return err
	}
	if response.Status.Code != 1 {
		return errors.New(response.Status.Message)
	}
	return nil
}

func (self *DNSPod) RecordInfo(domain, recordID string) (interface{}, error) {
	cli := ghttp.NewClient()
	params := self.PublicParam()
	params["domain"] = domain
	params["record_id"] = recordID
	resp, err := cli.Post(fmt.Sprintf("%s/Record.Info", DNSPodApiUrl), params)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	response := DNSPodResponse{}
	err = gjson.New(resp.ReadAllString()).ToStruct(&response)
	return response, err
}

func DefaultDnsPod(id, token string) *DNSPod {
	return &DNSPod{
		ID:           id,
		Token:        token,
		Format:       "json",
		Lang:         "en",
		ErrorOnEmpty: "yes",
	}
}
