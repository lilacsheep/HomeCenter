package models

import (
	"errors"
	"fmt"
	"strings"

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

func (v *Version) Next() (string, error) {
	v.latest += 1
	nv := fmt.Sprintf("%d.%d.%d", v.first, v.second, v.latest)
	return nv, nil
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

type ObjectInfo struct {
	DefaultModel
	Name        string `json:"name"`
	Key         string `json:"key"`
	Size        int64  `json:"size"`
	Hash        string `json:"hash"`
	RealPath    string `json:"real_path"`
	ContextType string `json:"context_type"`
	Version     string `json:"version"`
	gmeta.Meta  `orm:"table:object_table"`
}

func (s *ObjectInfo) GetVersion() (*Version, error) {
	return NewVersion(s.Version)
}
