package mallory

import (
	"errors"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"homeproxy/library/filedb"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	ErrShouldProxy = errors.New("should proxy")
)

type DomainErrorEvent struct {
	Domain string
	Error  string
}

func (self *DomainErrorEvent) DoEvent() error {
	params := g.Map{"domain": self.Domain, "error": self.Error}
	var data = ProxyRoleAnalysis{}
	err := filedb.DB.QueryOne(ProxyRoleAnalysisTable, &data, params)
	if err != nil {
		if err != filedb.ErrNoData {
			return err
		} else {
			data.Domain = self.Domain
			data.Times = 1
			data.Error = self.Error
			_, err := filedb.DB.Insert(ProxyRoleAnalysisTable, data)
			if err != nil {
				return err
			}
		}
	} else {
		data.Times += 1
		err := filedb.DB.UpdateById(ProxyRoleAnalysisTable, data.ID, data)
		if err != nil {
			return err
		}
	}
	return nil
}

type closeWriter interface {
	CloseWrite() error
}

// Direct fetcher
type Direct struct {
	Tr *http.Transport
}

// Create and initialize
func NewDirect(shouldProxyTimeout time.Duration) *Direct {
	if shouldProxyTimeout == 0 {
		shouldProxyTimeout = 200 * time.Millisecond
	}
	tr := http.DefaultTransport.(*http.Transport)
	tr.Dial = (&net.Dialer{
		Timeout: shouldProxyTimeout,
	}).Dial
	return &Direct{Tr: tr}
}

// Data flow:
//  1. Receive request R1 from client
//  2. Re-post request R1 to remote server(the one client want to connect)
//  3. Receive response P1 from remote server
//  4. Send response P1 to client
func (self *Direct) ServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	if r.Method == "CONNECT" {
		glog.Println("this function can not handle CONNECT method")
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()

	// Client.Do is different from DefaultTransport.RoundTrip ...
	// Client.Do will change some part of request as a new request of the server.
	// The underlying RoundTrip never changes anything of the request.
	resp, err := self.Tr.RoundTrip(r)
	if err != nil {
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			glog.Errorf("RoundTrip: %s, reproxy...\n", err.Error())
			err = ErrShouldProxy
			return
		}
		glog.Errorf("RoundTrip: %s\n", err.Error())
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// please prepare header first and write them
	CopyHeader(w, resp)
	w.WriteHeader(resp.StatusCode)

	n, err := io.Copy(w, resp.Body)
	if err != nil {
		glog.Printf("Copy: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d := BeautifyDuration(time.Since(start))
	ndtos := BeautifySize(n)
	glog.Infof("RESPONSE %s %s in %s <-%s", r.URL.Host, resp.Status, d, ndtos)
	return
}

// Data flow:
//  1. Receive CONNECT request from the client
//  2. Dial the remote server(the one client want to conenct)
//  3. Send 200 OK to client if the connection is established
//  4. Exchange data between client and server
func (self *Direct) Connect(w http.ResponseWriter, r *http.Request) (err error) {
	if r.Method != "CONNECT" {
		glog.Println("this function can only handle CONNECT method")
		http.Error(w, r.Method, http.StatusMethodNotAllowed)
		return
	}
	start := time.Now()

	// Use Hijacker to get the underlying connection
	hij, ok := w.(http.Hijacker)
	if !ok {
		s := "Server does not support Hijacker"
		http.Error(w, s, http.StatusInternalServerError)
		return
	}

	// connect the remote client directly
	dst, err := self.Tr.Dial("tcp", r.URL.Host)
	if err != nil {
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			err = ErrShouldProxy
			glog.Infof("RoundTrip: %s, reproxy...", err.Error())
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	src, _, err := hij.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer src.Close()

	// Once connected successfully, return OK
	src.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	// Proxy is no need to know anything, just exchange data between the client
	// the the remote server.
	copyAndWait := func(dst, src net.Conn, c chan int64) {
		n, err := io.Copy(dst, src)
		if err != nil {
			glog.Debugf("Copy: %s\n", err.Error())
			// FIXME: how to report error to dst ?
		}
		if tcpConn, ok := dst.(closeWriter); ok {
			tcpConn.CloseWrite()
		}
		c <- n
	}

	// client to remote
	stod := make(chan int64)
	go copyAndWait(dst, src, stod)

	// remote to client
	dtos := make(chan int64)
	go copyAndWait(src, dst, dtos)

	var nstod, ndtos int64
	for i := 0; i < 2; {
		select {
		case nstod = <-stod:
			i++
		case ndtos = <-dtos:
			i++
		}
	}
	d := BeautifyDuration(time.Since(start))
	glog.Infof("CLOSE %s after %s ->%s <-%s",
		r.URL.Host, d, BeautifySize(nstod), BeautifySize(ndtos))
	return
}
