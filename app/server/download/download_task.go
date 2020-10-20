package download

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/grpool"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	"os"
	"sync"
	"time"
)

type Task struct {
	ID          string `json:"id"`
	Url         string `json:"url"`
	FileName    string `json:"file_name"`
	Multi       bool   `json:"multi"` // 多线程下载
	helper      *ChunkHelper
	Status      int            `json:"status"` // 1 等待中 2 下载中 3 完成 99 失败
	Chunks      map[int]*Chunk `json:"chunks"`
	TotalSize   int            `json:"total_size"`
	Progress    string         `json:"progress"`
	Speed       int            `json:"speed"`
	doneSize    int
	DoneChan    chan bool `json:"-"`
	SizeChan    chan int  `json:"-"`
	ChunkSize   int       `json:"chunk_size"`
	fileHandler *os.File
}

func (self *Task) chunk() {
	for {
		start, end := self.helper.Next()
		if start == 0 && end == 0 {
			break
		} else {
			self.Chunks[start] = &Chunk{Start: start, End: end}
		}
	}
}

func (self *Task) singleDownload() {
	cli := ghttp.NewClient()
	if response, err := cli.Get(self.Url); err != nil {
		glog.Errorf("init url err: %s", err.Error())
	} else {
		defer response.Body.Close()
		var buf [512]byte
		var end = false

		for !end {
			if nr, err := response.Body.Read(buf[:]); err != nil {
				glog.Error(err.Error())
				end = true
				self.DoneChan <- end
			} else {
				switch true {
				case nr < 0:
					end = true
					self.DoneChan <- end
					break
				case nr == 0:
					end = true
					self.DoneChan <- end
					break
				case nr > 0:
					self.fileHandler.Write(buf[0:nr])
					self.SizeChan <- nr
				}
			}
		}
	}
}

func (self *Task) sync() {
	done := 0
	sync_ := false
	go func() {
		for !sync_ {
			select {
			case <-self.DoneChan:
				done++
			case n := <-self.SizeChan:
				self.doneSize += n
			default:
				sync_ = (self.doneSize == self.TotalSize) && (self.ChunkSize == done)
			}
		}

	}()
	go func() {
		history := 0
		for !sync_ {
			p := (float64(self.doneSize) / float64(self.TotalSize)) * 100
			self.Progress = fmt.Sprintf("%.2f", p)
			self.Speed = self.doneSize - history
			history = self.doneSize
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second)
		p := (float64(self.doneSize) / float64(self.TotalSize)) * 100
		self.Progress = fmt.Sprintf("%.2f", p)
		self.Speed = 0
		self.Status = 3
	}()
}

func (self *Task) Do() {
	self.fileHandler, _ = gfile.OpenFile(self.FileName, os.O_CREATE|os.O_RDWR, 0755)
	defer self.fileHandler.Close()

	if self.Multi {
		self.chunk()
		self.ChunkSize = len(self.Chunks)
		pool := grpool.New(4)
		wg := sync.WaitGroup{}
		d := gmap.NewTreeMap(gutil.ComparatorInt, true)

		for k, v := range self.Chunks {
			if (v.End - v.Start) != v.DoneSize {
				d.Set(k, v)
			}
		}

		d.IteratorAsc(func(key, value interface{}) bool {
			wg.Add(1)
			pool.Add(func() {
				value.(*Chunk).Download(self.Url, self.fileHandler, self.SizeChan, self.DoneChan)
				wg.Done()
			})
			return true
		})
		self.sync()
		wg.Wait()
	} else {
		self.ChunkSize = 1
		self.sync()
		self.singleDownload()
	}

}

func NewDownLoadTask(url string, chunkSize int) (*Task, error) {
	multi := true
	cli := ghttp.NewClient()
	response, err := cli.Head(url)
	if err != nil {
		return nil, err
	}
	header := response.Header.Clone()
	defer response.Body.Close()
	name := gfile.Name(url)
	ext := gfile.ExtName(url)
	size := gconv.Int(header.Get("Content-Length"))
	if header.Get("Accept-Ranges") == "" {
		multi = false
	}
	chunks := make(map[int]*Chunk)
	return &Task{
		Url:       url,
		Multi:     multi,
		TotalSize: size,
		DoneChan:  make(chan bool, 1),
		SizeChan:  make(chan int, 99999),
		FileName:  fmt.Sprintf("%s.%s", name, ext),
		helper:    &ChunkHelper{Total: size, ChunkSize: chunkSize},
		Chunks:    chunks,
	}, nil
}
