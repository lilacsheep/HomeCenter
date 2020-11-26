package ddns

import "github.com/gogf/gf/frame/g"

type Provider interface {
	Version() (interface{}, error)
	PublicParam() g.Map
	DomainList(keyword ...string) (interface{}, error)
	RecordList(domain string, subDomain string, keyword ...string) (interface{}, error)
	// create record return record_id and error
	RecordCreate(domain, subDomain, ttl string, address Address) (string, error)
	RecordRemove(domain string, recordID string) (interface{}, error)
	RecordModify(domain, recordID, subDomain string, address Address) error
	RecordInfo(domain, recordID string) (interface{}, error)
}
