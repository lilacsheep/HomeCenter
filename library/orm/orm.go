package orm

import (
	"homeproxy/library/config"
	"path"

	"github.com/gogf/gf/database/gdb"
)

var (
	DB gdb.DB
)

func Init() error {
	dbPath := path.Join(config.Dbpath, "db.sqlite")
	gdb.SetConfig(gdb.Config{"default": gdb.ConfigGroup{
		gdb.ConfigNode {
            Host: dbPath,
			Type: "sqlite",
        },
	}})
	var err error

	DB, err = gdb.New()
	if err != nil {
		return err
	}
	return nil
}