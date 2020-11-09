package download

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"io"
	"net/http"
	"os"
	"time"
)

type Block struct {
	Buf  [512]byte
	Size int64
}

type Chunk struct {
	Start      int64 `json:"start"`
	End        int64 `json:"end"`
	DoneSize   int64 `json:"done_size"`   // 已经完成的字节数
	ExpectSize int64 `json:"expect_size"` // 分片预期下载的字节数
	Done       bool  `json:"done"`
	canceled   bool
}

func (self *Chunk) Cancel() {
	self.canceled = true
}

func (self *Chunk) header() string {
	var start int64
	if self.DoneSize > 0 {
		start = self.Start + self.DoneSize
	} else {
		start = self.Start
	}
	if self.End == 0 {
		return fmt.Sprintf("bytes=%d-%d", start, self.End)
	} else {
		return fmt.Sprintf("bytes=%d-%d", start, self.End)
	}
}

func (self *Chunk) checkResponse(response *ghttp.ClientResponse) bool {
	var check bool
	check = response.ContentLength == (self.ExpectSize - self.DoneSize)
	check = response.StatusCode == http.StatusPartialContent
	check = response.Header.Get("Accept-Ranges") == "bytes"
	return check
}

func (self *Chunk) download(url string, write chan Block, errChan chan error) {
	cli := ghttp.NewClient()
	header := self.header()
	cli.SetHeader("Range", header)
	response, err := cli.Get(url)
	if err != nil {
		errChan <- err
		return
	}
	defer response.Close()
	if self.checkResponse(response) {
		var buf [512]byte
		for !self.Done && !self.canceled {
			switch nr, _ := response.Body.Read(buf[:]); true {
			case nr > 0:
				write <- Block{Buf: buf, Size: gconv.Int64(nr)}
			case nr == 0:
				self.Done = self.checkDone()
			case nr < 0:
				errChan <- io.EOF
				return
			}
		}
	}
}

func (self *Chunk) checkDone() bool {
	return self.DoneSize == self.ExpectSize
}

func (self *Chunk) Download(url string, file *os.File) {
	self.canceled = false
	var offset = self.Start + self.DoneSize
	writeChan := make(chan Block)
	errChan := make(chan error)

Loop:
	if !self.checkDone() {
		go self.download(url, writeChan, errChan)
		for !self.Done && !self.canceled {
			select {
			case block := <-writeChan:
				file.WriteAt(block.Buf[0:block.Size], gconv.Int64(offset))
				offset += gconv.Int64(block.Size)
				self.DoneSize += gconv.Int64(block.Size)
			case <-errChan:
				if self.checkDone() {
					self.Done = true
				} else {
					glog.Warning("download error, restart 5s after")
					time.Sleep(5 * time.Second)
					goto Loop
				}
			default:
				self.Done = self.checkDone()
			}
		}
	}
}
