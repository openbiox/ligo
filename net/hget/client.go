package hget

import (
	"crypto/tls"
	"errors"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strings"
	"time"

	"github.com/openbiox/ligo/stringo"
)

var gCurCookies []*http.Cookie
var gCurCookieJar *cookiejar.Jar

func NewHTTPClient(timeout int, proxy string) *http.Client {
	transPort := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(timeout) * time.Second,
		}).Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy != "" {
		urlproxy, _ := RandProxy(proxy)
		transPort.Proxy = http.ProxyURL(urlproxy)
	}
	return &http.Client{
		CheckRedirect: defaultCheckRedirect,
		Jar:           gCurCookieJar,
		Transport:     transPort,
	}
}

// RandProxy return a proxy from proxy string
func RandProxy(proxy string) (*neturl.URL, string) {
	if proxy == "" {
		return nil, ""
	}
	proxyPool := []string{}
	if strings.Contains(proxy, ";") {
		proxyPool = stringo.StrSplit(proxy, ";", 1000000)
	} else {
		proxyPool = append(proxyPool, proxy)
	}
	i := rand.Int63n(int64(len(proxyPool) - 0))
	urli := neturl.URL{}
	urlproxy, _ := urli.Parse(proxyPool[i])
	return urlproxy, proxyPool[i]
}

func defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 20 {
		return errors.New("stopped after 20 redirects")
	}
	return nil
}

func init() {
	gCurCookies = nil
	//var err error;
	gCurCookieJar, _ = cookiejar.New(nil)
}
