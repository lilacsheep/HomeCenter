package download

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	url2 "net/url"
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

type Task struct {
	ID          int              `json:"id" storm:"id,increment"`
	Url         string           `json:"url" storm:"unique"`
	FileName    string           `json:"file_name"`
	FilePath    string           `json:"file_path"`
	Multi       bool             `json:"multi"`  // 多线程下载
	Status      uint             `json:"status"` // 1 暂停 2 等待中 3 下载中 4 完成 99 失败
	Chunks      map[int64]*Chunk `json:"chunks"`
	TotalSize   int64            `json:"total_size"`
	Progress    string           `json:"progress"`
	Speed       int64            `json:"speed"`
	DoneSize    int64            `json:"done_size"`
	MD5         string           `json:"md5"`
	ThreadNum   int64            `json:"thread_num"`
	fileHandler *os.File
	Canceled    bool `json:"canceled"`
}

func (self *Task) chunk() {
	mod := self.TotalSize % self.ThreadNum
	blockSize := (self.TotalSize - mod) / self.ThreadNum
	var i int64
	for i = 1; i < self.ThreadNum+1; i++ {
		if i == self.ThreadNum {
			self.Chunks[i] = &Chunk{Start: blockSize * (i - 1), End: self.TotalSize, ExpectSize: blockSize + mod}
		} else {
			self.Chunks[i] = &Chunk{Start: blockSize * (i - 1), End: blockSize*i - 1, ExpectSize: blockSize}
		}
	}
}

func (self *Task) singleDownload() {
	cli := ghttp.NewClient()
	if response, err := cli.Get(self.Url); err != nil {
		glog.Errorf("init url err: %s", err.Error())
	} else {
		self.Status = 3
		defer response.Body.Close()
		var buf [512]byte
		var end = false

		for !end && !self.Canceled {
			if nr, err := response.Body.Read(buf[:]); err != nil {
				glog.Error(err.Error())
				end = true
			} else {
				switch true {
				case nr < 0:
					end = true
					break
				case nr == 0:
					end = true
					break
				case nr > 0:
					self.fileHandler.Write(buf[0:nr])
					self.DoneSize += gconv.Int64(nr)
				}
			}
		}
	}
}

func (self *Task) Cancel() {
	self.Canceled = true

	if self.Multi {
		d := gmap.NewTreeMap(gutil.ComparatorInt, true)

		for k, v := range self.Chunks {
			d.Set(k, v)
		}

		d.IteratorAsc(func(key, value interface{}) bool {
			value.(*Chunk).Cancel()
			return true
		})
	}
	self.Status = 1
}

func (self *Task) checkChunk() (int64, int64) {
	var (
		nowSize int64
		done    int64
	)
	for _, chunk := range self.Chunks {
		nowSize += chunk.DoneSize
		if chunk.Done {
			done += 1
		}
	}
	return nowSize, done
}

func (self *Task) sync() {
	go func() {
		var (
			history int64
			done    int64
		)
		for done != self.ThreadNum && !self.Canceled {
			if self.Multi {
				self.DoneSize, done = self.checkChunk()
			}
			self.Progress = fmt.Sprintf("%.2f", (float64(self.DoneSize)/float64(self.TotalSize))*100)
			self.Speed = self.DoneSize - history
			history = self.DoneSize
			time.Sleep(time.Second)
		}
		if self.Multi {
			self.DoneSize, done = self.checkChunk()
		}
		self.Progress = fmt.Sprintf("%.2f", (float64(self.DoneSize)/float64(self.TotalSize))*100)
		self.Speed = 0
		if !self.Canceled {
			self.Status = 4
		}
	}()
}

func (self *Task) Init() {
	if self.Multi {
		if len(self.Chunks) == 0 {
			self.chunk()
		}
	} else {
		self.ThreadNum = 1
	}
}

func (self *Task) Start() {
	filename := gfile.Join(self.FilePath, self.FileName)
	self.fileHandler, _ = gfile.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0755)
	self.Canceled = false
	if self.Multi {
		wg := sync.WaitGroup{}
		d := gmap.NewTreeMap(gutil.ComparatorInt, true)
		for k, v := range self.Chunks {
			d.Set(k, v)
		}
		self.Status = 3
		d.IteratorAsc(func(key, value interface{}) bool {
			wg.Add(1)
			go func() {
				chunk := value.(*Chunk)
				chunk.Download(self.Url, self.fileHandler)
				wg.Done()
			}()
			return true
		})
		time.Sleep(time.Millisecond * 500)
		self.sync()
		wg.Wait()
	} else {
		self.sync()
		self.singleDownload()
	}
	_ = self.fileHandler.Close()

	if !self.Canceled {
		h := sha256.New()
		f, _ := gfile.Open(gfile.Join(self.FilePath, self.FileName))
		defer f.Close()
		_, err := io.Copy(h, f)
		if err != nil {
			glog.Errorf("check sum error: %s", err.Error())
		} else {
			self.MD5 = hex.EncodeToString(h.Sum(nil))
		}
	}
}

func NewDownLoadTask(url string, threadNum int64, path string) (*Task, error) {
	urlCode, _ := url2.QueryUnescape(url)
	multi := true
	cli := ghttp.NewClient()
	response, err := cli.Head(url)
	if err != nil {
		return nil, err
	}
	header := response.Header.Clone()
	defer response.Body.Close()

	name := gfile.Name(urlCode)
	ext := gfile.ExtName(urlCode)
	if header.Get("Accept-Ranges") == "" {
		multi = false
	}
	chunks := make(map[int64]*Chunk)
	return &Task{
		Url:       url,
		Multi:     multi,
		TotalSize: response.ContentLength,
		FileName:  fmt.Sprintf("%s.%s", name, ext),
		FilePath:  path,
		ThreadNum: threadNum,
		Chunks:    chunks,
		Status:    1,
	}, nil
}
