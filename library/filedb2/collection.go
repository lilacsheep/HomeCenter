package filedb2

import (
	"fmt"
	"homeproxy/library/config"
	"path/filepath"

	"github.com/asdine/storm/v3"
)

var (
	DB *storm.DB
)

func Init() {
	dbname := fmt.Sprintf("%s.db", config.Dbname)
	DB, _ = storm.Open(filepath.Join(config.Dbpath, dbname))
}