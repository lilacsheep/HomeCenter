package filedb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/guid"
	"github.com/gogf/gf/util/gutil"
	"io/ioutil"
	"path"
	"sync"
)

var (
	ErrCollectionExist    = errors.New("collection exist")
	ErrCollectionNotExist = errors.New("collection not exist")
	ErrNoData             = errors.New("no data")
)

type FileSource struct {
	FilePath string
	lock     *sync.RWMutex
}

func (self *FileSource) Read() ([]byte, error) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if gfile.Exists(self.FilePath) {
		return ioutil.ReadFile(self.FilePath)
	}
	return nil, nil
}

func (self *FileSource) Save(data []byte) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	return ioutil.WriteFile(self.FilePath, data, 0755)
}

type Collection struct {
	file *FileSource
	Name string
	data *gmap.TreeMap
}

func (self *Collection) Insert(data interface{}) {
	j := gjson.New(data, true)
	j.Set("id", guid.S())
	self.data.Set(j.Get("id"), j)
}

func (self *Collection) GetFirst(data interface{}) error {
	if self.data.Size() > 0 {
		v, _ := self.data.Search(self.data.Keys()[0])
		return v.(*gjson.Json).ToStruct(data)
	}
	return ErrNoData
}

func (self *Collection) GetById(id string, data interface{}) error {
	return self.data.Get(id).(*gjson.Json).ToStruct(data)
}

func (self *Collection) RemoveById(id string) {
	self.data.Remove(id)
}

func (self *Collection) GetAndRemove(id string, data interface{}) error {
	return self.data.Remove(id).(*gjson.Json).ToStruct(data)
}

func (self *Collection) All(data interface{}) error {
	temp := gjson.New(nil)
	index := 0
	self.data.Iterator(func(key, value interface{}) bool {
		for k, v := range value.(*gjson.Json).ToMap() {
			temp.Set(fmt.Sprintf("%d.%s", index, k), v)
		}
		index++
		return true
	})
	return json.Unmarshal([]byte(temp.MustToJsonString()), data)
}

func (self *Collection) Search(params g.Map, data interface{}) error {
	temp := gjson.New(nil)
	index := 0
	self.data.Iterator(func(key, value interface{}) bool {
		same := true
		t := value.(*gjson.Json)
		for k, v := range params {
			if t.Get(k) != v {
				same = false
			}
		}
		if same {
			for k, v := range t.ToMap() {
				temp.Set(fmt.Sprintf("%d.%s", index, k), v)
			}
			index++
		}
		return true
	})
	return json.Unmarshal([]byte(temp.MustToJsonString()), data)
}

func (self *Collection) Paging(offset, limit int, data interface{}) (err error) {
	index := 0
	d := gjson.New(nil)
	if offset >= self.data.Size() {
		var keys []interface{}

		if len(self.data.Keys()[offset:]) >= limit {
			keys = self.data.Keys()[offset : offset+limit]
		} else {
			keys = self.data.Keys()[offset : offset+len(self.data.Keys()[offset:])]
		}

		for _, key := range keys {
			if v, ok := self.data.Search(key); ok {
				d.Set(fmt.Sprintf("%d", index), v.(*gjson.Json).ToMap())
				index++
			}
		}
	}
	return json.Unmarshal([]byte(d.MustToJsonString()), data)
}

func (self *Collection) UpdateById(id string, data interface{}) {
	g_ := gjson.New(data, true)
	if v, ok := self.data.Search(id); ok {
		ov := v.(*gjson.Json)
		for k, value := range g_.ToMap() {
			ov.Set(k, value)
		}
		self.data.Set(id, ov)
	}
}

func (self *Collection) Load() {
	data, err := self.file.Read()
	if err != nil {
		glog.Errorf("load collection %s error", self.Name, err)
	} else {
		if data == nil {
			glog.Warning("load database mapping is null")
		} else {
			d := gjson.New(data, true)
			for k, v := range d.ToMap() {
				self.data.Set(k, gjson.New(v, true))
			}
		}
	}
}

func (self *Collection) Dump() {
	t := gjson.New(nil, false)
	self.data.Iterator(func(key, value interface{}) bool {
		t.Set(key.(string), value.(*gjson.Json).ToMap())
		return true
	})
	data := t.MustToJson()
	if err := self.file.Save(data); err != nil {
		glog.Errorf("dump collection error: %s", err.Error())
	}
}

type Database struct {
	collections *gmap.StrAnyMap
	file        *FileSource
	Path        string
	Name        string
}

func (self *Database) Collection(key string) (*Collection, error) {
	if v, ok := self.collections.Search(key); !ok {
		return nil, ErrCollectionNotExist
	} else {
		return v.(*Collection), nil
	}
}

func (self *Database) NewCollections(name string) (err error) {
	if _, ok := self.collections.Search(name); !ok {
		self.collections.Set(name, &Collection{
			Name: name,
			data: gmap.NewTreeMap(gutil.ComparatorString, true),
			file: &FileSource{
				FilePath: path.Join(self.Path, fmt.Sprintf("%s.json", name)),
				lock:     &sync.RWMutex{},
			},
		})
	} else {
		return ErrCollectionExist
	}
	return
}

func (self *Database) Load() {
	NameMap := garray.NewStrArray(true)
	data, err := self.file.Read()
	if err != nil {
		glog.Errorf("load database mapping error: %s", err.Error())
	} else {
		if data == nil {
			glog.Warning("load database mapping is null")
		} else {
			if err := NameMap.UnmarshalJSON(data); err != nil {
				glog.Errorf("load database mapping error: %s", err.Error())
			} else {
				NameMap.Iterator(func(k int, v string) bool {
					if err := self.NewCollections(v); err != nil {
						glog.Errorf("set collections error: %s", err.Error())
					}
					return true
				})
				self.collections.Iterator(func(k string, v interface{}) bool {
					v.(*Collection).Load()
					return true
				})
			}
		}
	}
}

func (self *Database) Dump() {
	if d, err := json.Marshal(self.collections.Keys()); err == nil {
		if err := self.file.Save(d); err != nil {
			glog.Errorf("dump collection mapping error: %s", err.Error())
		}
		self.collections.Iterator(func(k string, v interface{}) bool {
			collection := v.(*Collection)
			collection.Dump()
			return true
		})
	} else {
		glog.Errorf("dump database file error: %s", err)
	}
}

func NewDatabase(name, path string) *Database {
	database := &Database{
		collections: gmap.NewStrAnyMap(true),
		file:        &FileSource{FilePath: gfile.Join(path, fmt.Sprintf("%s_collection.json", name)), lock: &sync.RWMutex{}},
		Name:        name,
		Path:        path,
	}
	database.Load()
	return database
}
