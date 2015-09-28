package main

import (
	"os"
	"strings"

	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/koding/websocketproxy"
)

var HubURL = os.Getenv("HUB_URL")
var SparkURL = os.Getenv("SPARK_URL")

func getFQDN(host string, port string) string {
	return host + ":" + port
}

func newHubProxy() *httputil.ReverseProxy {
	return newReverseProxyUrlRewrite(HubURL, func(oldUrl string, url *url.URL, c []*http.Cookie) {
		_, username := getPortUsername(c)
		log.Debugf("User: %s, proxy Hub: %s -> %s", username, oldUrl, url.String())
	})
}

func newSparkMasterProxy() *httputil.ReverseProxy {
	return newReverseProxyUrlRewrite(SparkURL, func(oldUrl string, url *url.URL, c []*http.Cookie) {
		_, username := getPortUsername(c)
		url.Path = strings.Replace(url.Path, "/api/v1/cluster", "/json/", 1)
		log.Debugf("User: %s, proxy Spark: %s -> %s", username, oldUrl, url.String())
	})
}

func newZeppelinProxy() *httputil.ReverseProxy {
	return newReverseProxyUrlRewrite("http://localhost",
		func(oldUrl string, url *url.URL, cookies []*http.Cookie) {

			if strings.Contains(url.Path, "/zeppelin/") {
				url.Path = strings.Replace(url.Path, "/zeppelin/", "/", 1)
			} else {
				url.Path = strings.Replace(url.Path, "/zeppelin", "/", 1)
			}

			port, username := getPortUsername(cookies)
			host := userHosts[username]
			url.Host = getFQDN(host, port)
			log.Debugf("User: %s, proxy Zeppelin: %s -> %s", username, oldUrl, url.String())
		})
}

type UrlRewriteWithCookies func(oldUrl string, url *url.URL, cookies []*http.Cookie)

func newReverseProxyUrlRewrite(urlStr string, rewrite UrlRewriteWithCookies) *httputil.ReverseProxy {
	url := parseUrl(urlStr)
	rp := httputil.NewSingleHostReverseProxy(url)
	oldDirector := rp.Director
	rp.Director = func(r *http.Request) {
		oldDirector(r)
		oldUrl := r.URL.String()
		r.Host = url.Host
		rewrite(oldUrl, r.URL, r.Cookies())
	}
	return rp
}

func newWebsocketProxy() *websocketproxy.WebsocketProxy {
	wp := websocketproxy.NewProxy(parseUrl("ws://localhost"))

	wp.Backend = func(r *http.Request) *url.URL {
		u := *r.URL
		u.Fragment = r.URL.Fragment
		u.Path = r.URL.Path
		u.RawQuery = r.URL.RawQuery

		port, username := getPortUsername(r.Cookies())
		log.Tracef("user:%v port:%v found in %+v", username, port, r.Cookies())
		host := userHosts[username]
		u.Host = getFQDN(host, port)
		u.Scheme = "ws"
		log.Debugf("User: %s, proxy websocket: %s -> %s", username, r.URL.String(), u.String())
		return &u
	}

	return wp
}

func parseUrl(urlStr string) *url.URL {
	url, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal("URL " + urlStr + " failed to parse")
	}
	return url
}

func getPortUsername(cookies []*http.Cookie) (string, string) {
	port := "65535"
	username := "guest"
	for _, cookie := range cookies {
		if cookie.Name == "port" {
			port = cookie.Value
		} else if cookie.Name == "username" {
			username = cookie.Value
		}
	}
	return port, username
}
