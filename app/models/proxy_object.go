package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gmeta"
)

type Version struct {
	first  int
	second int
	latest int
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.first, v.second, v.latest)
}

func (v Version) Ints() (int, int, int) {
	return v.first, v.second, v.latest
}

func (v Version) Gte(ver Version) bool {
	first, second, latest := ver.Ints()
	if v.first == first && v.second == second {
		return v.latest >= latest
	}
	if v.first == first {
		return v.second >= second
	}
	return v.first >= first
}

func (v Version) Gt(ver Version) bool {
	first, second, latest := ver.Ints()
	if v.first == first && v.second == second {
		return v.latest > latest
	}
	if v.first == first {
		return v.second > second
	}
	return v.first > first
}

func (v Version) Equal(ver Version) bool {
	first, second, latest := ver.Ints()
	return v.first == first && v.second == second && v.latest == latest
}

func (v *Version) Next() Version {
	return Version{v.first, v.second, v.latest + 1}
}

func NewVersion(data string) (*Version, error) {
	if data == "" {
		return &Version{first: 1, second: 0, latest: 0}, nil
	}
	t := strings.Split(data, ".")
	if len(t) != 3 {
		return nil, errors.New("error format")
	}
	return &Version{first: gconv.Int(t[0]), second: gconv.Int(t[1]), latest: gconv.Int(t[2])}, nil
}

type ObjectToken struct {
	DefaultModel
	Name       string `json:"name"`
	Effective  int    `json:"effective"`
	SecretKey  string `json:"secret_key"`
	Upload     bool   `json:"upload"`
	Download   bool   `json:"download"`
	Delete     bool   `json:"delete"`
	List       bool   `json:"list"`
	gmeta.Meta `orm:"table:object_token"`
}

type Bucket struct {
	DefaultModel
	Name       string `json:"name"`
	Public     bool   `json:"public"`
	Referer    bool   `json:"referer"`
	RefererUrl string `json:"referer_url"`
	gmeta.Meta `orm:"table:object_bucket"`
}

type ObjectInfo struct {
	DefaultModel
	Name        string `json:"name"`
	Key         string `json:"key"`
	Size        int64  `json:"size"`
	Bucket      int    `json:"bucket"`
	Hash        string `json:"hash"`
	RealPath    string `json:"real_path"`
	ContextType string `json:"context_type"`
	Version     string `json:"version"`
	gmeta.Meta  `orm:"table:objects"`
}

func (s *ObjectInfo) GetVersion() (*Version, error) {
	return NewVersion(s.Version)
}

func (s *ObjectInfo) CopyNewRecord(bucket int, name, key string) g.Map {
	return g.Map{
		"name":         name,
		"key":          key,
		"bucket":       bucket,
		"hash":         s.Hash,
		"size":         s.Size,
		"real_path":    s.RealPath,
		"context_type": s.ContextType,
		"verison":      s.Version,
	}
}

// 分享文件
type ShareFileList struct {
	DefaultModel
	ID       int       `json:"id"`
	Vkey     string    `json:"vkey"`
	ObjectId int       `json:"object_id"`
	CreateAt time.Time `json:"create_at"`
}
