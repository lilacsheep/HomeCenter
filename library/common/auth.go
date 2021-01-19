package common

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"unsafe"
)

const (
	BasicAuthName         = "Basic"
	ProxyAuthorizationKey = "Proxy-Authorization"
	ProxyAuthenticateKey  = "Proxy-Authenticate"
)

// AuthenticationFunc Authentication interface is implemented
type AuthenticationFunc func(http.ResponseWriter, *http.Request) bool

// Auth authentication processing
func (f AuthenticationFunc) Auth(w http.ResponseWriter, r *http.Request) bool {
	return f(w, r)
}

// Authentication proxy authentication
type Authentication interface {
	Auth(http.ResponseWriter, *http.Request) bool
}

// BasicAuth HTTP Basic authentication for Header Proxy-Authorization
func BasicAuth(username, password string) Authentication {
	return BasicAuthFunc(func(u, p string) bool {
		return username == u && password == p
	})
}

// BasicAuthFunc HTTP Basic authentication for Header Proxy-Authorization
func BasicAuthFunc(f func(username, password string) bool) Authentication {
	return AuthenticationFunc(func(w http.ResponseWriter, r *http.Request) bool {
		if u, p, ok := parseBasicAuth(r.Header.Get(ProxyAuthorizationKey)); ok && f(u, p) {
			return true
		}
		w.Header().Set(ProxyAuthenticateKey, BasicAuthName)
		http.Error(w, http.StatusText(http.StatusProxyAuthRequired), http.StatusProxyAuthRequired)
		return false
	})
}

// parseBasicAuth parses an HTTP Basic Authentication string.
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = BasicAuthName + " "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := *(*string)(unsafe.Pointer(&c))
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

// basicAuth HTTP Basic Authentication string.
func basicAuth(u *url.Userinfo) (base string) {
	const prefix = BasicAuthName + " "
	s := u.String()
	base = base64.StdEncoding.EncodeToString(*(*[]byte)(unsafe.Pointer(&s)))
	return prefix + base
}