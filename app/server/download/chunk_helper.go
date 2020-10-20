package download

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"os"
)

type Chunk struct {
	Start    int `json:"start"`
	End      int `json:"end"`
	DoneSize int `json:"done_size"`
}

func (self *Chunk) Download(url string, file *os.File, sizeChan chan int, done chan bool) {
	var offset = self.Start

	cli := ghttp.NewClient()
	cli.SetHeader("Range", fmt.Sprintf("bytes=%d-%d", self.Start, self.End))
	if response, err := cli.Get(url); err != nil {
		glog.Errorf("init url err: %s", err.Error())
	} else {
		defer response.Body.Close()

		var buf [512]byte
		var end = false

		for !end {
			switch nr, _ := response.Body.Read(buf[:]); true {
			case nr < 0:
				end = true
				done <- end
				break
			case nr == 0:
				end = true
				done <- end
				break
			case nr > 0:
				file.WriteAt(buf[0:nr], gconv.Int64(offset))
				offset += nr
				self.DoneSize += nr
				sizeChan <- nr
			}
		}
	}
}

type ChunkHelper struct {
	size      int
	Total     int
	ChunkSize int
	end       bool
}

func (self *ChunkHelper) Next() (int, int) {
	if self.end {
		return 0, 0
	}
	next := self.size + self.ChunkSize
	if next >= self.Total {
		self.end = true
		return self.size, self.Total
	} else {
		start, end := self.size, self.size+self.ChunkSize
		self.size += self.ChunkSize + 1
		return start, end
	}
}
