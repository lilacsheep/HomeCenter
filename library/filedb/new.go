package filedb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
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
	ErrUnique             = errors.New("unique error")
	ErrUniqueNotExist     = errors.New("unique field exist but index cache not exist")
)

type FileSource struct {
	FilePath string
	lock     *sync.RWMutex
}

func (self *FileSource) Size() int64 {
	return gfile.Size(self.FilePath)
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
	file        *FileSource
	Name        string
	data        *gmap.TreeMap
	Settings    *CollectionSettings
	UniqueIndex *gset.Set
	lock        *sync.RWMutex
}

func (self *Collection) checkUnique(value *gjson.Json) (bool, error) {
	if self.Settings.Unique != "" {
		if self.UniqueIndex != nil {
			if value.Contains(self.Settings.Unique) {
				return self.UniqueIndex.Contains(value.Get(self.Settings.Unique)), nil
			}
			return false, nil
		} else {
			return false, ErrUniqueNotExist
		}
	}
	return false, nil
}

func (self *Collection) insert(data interface{}) (id string, err error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	j := gjson.New(data)
	if found, err := self.checkUnique(j); err != nil {
		return "", err
	} else {
		if !found {
			id = guid.S()
			j.Set("id", id)
			self.data.Set(j.Get("id"), j)
			if self.UniqueIndex != nil {
				self.UniqueIndex.Add(j.Get(self.Settings.Unique))
			}
			self.dump(false)
			return id, nil
		} else {
			return "", ErrUnique
		}
	}
}

func (self *Collection) Insert(data interface{}) (string, error) {
	if self.Settings.MaxRecord != 0 && self.data.Size() < self.Settings.MaxRecord {
		return self.insert(data)
	} else if self.Settings.MaxRecord == 0 {
		return self.insert(data)
	} else {
		self.RemoveFirst()
		return self.insert(data)
	}
}

func (self *Collection) GetFirst(data interface{}) error {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if self.data.Size() > 0 {
		v, _ := self.data.Search(self.data.Keys()[0])
		return v.(*gjson.Json).ToStruct(data)
	}
	return ErrNoData
}

func (self *Collection) GetById(id string, data interface{}) error {
	self.lock.RLock()
	defer self.lock.RUnlock()
	v := self.data.Get(id)
	if v == nil {
		return ErrNoData
	}
	return v.(*gjson.Json).ToStruct(data)
}

func (self *Collection) RemoveFirst() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.data.Remove(self.data.Keys()[0])
	self.dump(false)
}

func (self *Collection) RemoveById(id string) {
	self.lock.Lock()
	defer self.lock.Unlock()
	v := self.data.Remove(id)
	if v != nil && self.UniqueIndex != nil {
		self.UniqueIndex.Remove(v.(*gjson.Json).Get(self.Settings.Unique))
	}
	self.dump(false)
}

func (self *Collection) GetAndRemove(id string, data interface{}) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	value := self.data.Remove(id)
	if value != nil {
		t := value.(*gjson.Json)
		if self.UniqueIndex != nil {
			self.UniqueIndex.Remove(t.Get(self.Settings.Unique))
		}
		self.dump(false)
		return t.ToStruct(data)
	} else {
		return ErrNoData
	}
}

func (self *Collection) All(data interface{}) error {
	self.lock.RLock()
	defer self.lock.RUnlock()
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
	self.lock.RLock()
	defer self.lock.RUnlock()
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
	self.lock.RLock()
	defer self.lock.RUnlock()
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
				t := v.(*gjson.Json)
				for k, v := range t.ToMap() {
					d.Set(fmt.Sprintf("%d.%s", index, k), v)
				}
				index++
			}
		}
	}
	return json.Unmarshal([]byte(d.MustToJsonString()), data)
}

func (self *Collection) UpdateById(id string, data interface{}) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	g_ := gjson.New(data, true)
	if found, err := self.checkUnique(g_); err != nil {
		return err
	} else {
		if found {
			return ErrUnique
		} else {
			if v, ok := self.data.Search(id); ok {
				ov := v.(*gjson.Json)
				for k, value := range g_.ToMap() {
					if k == self.Settings.Unique {
						self.UniqueIndex.Remove(ov.Get(k))
						self.UniqueIndex.Add(value)
					}
					ov.Set(k, value)
				}
				self.data.Set(id, ov)
				self.dump(false)
			}
		}
	}
	return nil
}

func (self *Collection) Load() {
	data, err := self.file.Read()
	if err != nil {
		glog.Errorf("load collection %s error", self.Name, err)
	} else {
		if data == nil {
			if self.Settings.AutoDump {
				glog.Warningf("load collection %s mapping is null", self.Name)
			}
		} else {
			d := gjson.New(data, true)
			if self.Settings.Unique != "" {
				if self.UniqueIndex != nil {
					for k, v := range d.ToMap() {
						value := gjson.New(v, true)
						found, err := self.checkUnique(value)
						if err != nil {
							glog.Errorf("%s data error %s data: %s", self.Name, err, value.MustToJson())
						} else {
							if found {
								glog.Errorf("%s data error %s", self.Name, ErrUnique)
							} else {
								self.data.Set(k, value)
								self.UniqueIndex.Add(value.Get(self.Settings.Unique))
							}
						}
					}
				} else {
					glog.Errorf("index set not exist")
				}
			} else {
				for k, v := range d.ToMap() {
					value := gjson.New(v, true)
					self.data.Set(k, value)
				}
			}
		}
	}
}

func (self *Collection) dump(force bool) {
	if self.Settings.AutoDump || force {
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
}

type Database struct {
	Cache       *gcache.Cache // key/value memory cache,
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

func (self *Database) NewCollections(name string, settings *CollectionSettings) (err error) {
	if _, ok := self.collections.Search(name); !ok {
		if settings == nil {
			settings = DefaultCollectionSettings()
		}
		var indexSet *gset.Set
		if settings.Unique != "" {
			indexSet = gset.New(true)
		}
		self.collections.Set(name, &Collection{
			Name:        name,
			data:        gmap.NewTreeMap(gutil.ComparatorString, true),
			Settings:    settings,
			UniqueIndex: indexSet,
			lock:        &sync.RWMutex{},
			file: &FileSource{
				FilePath: path.Join(self.Path, fmt.Sprintf("%s.json", name)),
				lock:     &sync.RWMutex{},
			},
		})
		self.Dump()
	} else {
		return ErrCollectionExist
	}
	return
}

func (self *Database) Load() {
	NameMap := gmap.NewStrAnyMap(true)
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
				NameMap.Iterator(func(k string, v interface{}) bool {
					settings := CollectionSettings{}
					if err := gjson.New(v).ToStruct(&settings); err == nil {
						if err := self.NewCollections(k, &settings); err != nil {
							glog.Errorf("set collections error: %s", err.Error())
						}
					} else {
						glog.Errorf("format error: %s", err.Error())
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
	NameMap := gmap.NewStrAnyMap()
	self.collections.Iterator(func(k string, v interface{}) bool {
		collection := v.(*Collection)
		NameMap.Set(collection.Name, collection.Settings)
		return true
	})

	if d, err := NameMap.MarshalJSON(); err == nil {
		if err := self.file.Save(d); err != nil {
			glog.Errorf("save collection mapping error: %s", err.Error())
		}
	} else {
		glog.Errorf("dump database file error: %s", err)
	}
}

func (self *Database) Insert(collectionName string, data interface{}) (id string, err error) {
	var collection *Collection
	collection, err = self.Collection(collectionName)
	if err != nil {
		return "", err
	}
	return collection.Insert(data)
}

func (self *Database) UpdateById(collectionName string, id string, data interface{}) (err error) {
	var collection *Collection
	collection, err = self.Collection(collectionName)
	if err != nil {
		return err
	}
	return collection.UpdateById(id, data)
}

func (self *Database) GetById(collectionName string, id string, data interface{}) (err error) {
	var collection *Collection
	collection, err = self.Collection(collectionName)
	if err != nil {
		return err
	}
	return collection.GetById(id, data)
}

func (self *Database) RemoveByID(collectionName string, id string) (err error) {
	var collection *Collection
	collection, err = self.Collection(collectionName)
	if err != nil {
		return err
	}
	collection.RemoveById(id)
	return nil
}

func (self *Database) QueryAll(collectionName string, data interface{}) (err error) {
	var collection *Collection
	collection, err = self.Collection(collectionName)
	if err != nil {
		return err
	}
	return collection.All(data)
}

func NewDatabase(name, path string) *Database {
	database := &Database{
		Cache:       gcache.New(),
		collections: gmap.NewStrAnyMap(true),
		file:        &FileSource{FilePath: gfile.Join(path, fmt.Sprintf("%s_collection.json", name)), lock: &sync.RWMutex{}},
		Name:        name,
		Path:        path,
	}
	database.Load()
	return database
}
