package cache

import (
	"homeproxy/library/config"

	"github.com/roseduan/rosedb"
)

var (
	DB *rosedb.RoseDB
)

func Init() error {
	var err error
	cfg := rosedb.DefaultConfig()
	cfg.DirPath = config.Dbpath
	DB, err = rosedb.Open(cfg)
	return err
}
