package utils

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func ProxyToService(targetBaseUrl string, pathPrefix string) http.HandlerFunc {
	target,_:=url.Parse(targetBaseUrl)

	proxy:=httputil.NewSingleHostReverseProxy(target)

	originalDirector:=proxy.Director

	proxy.Director=func(r *http.Request) {
		originalDirector(r)

		originalPath:=r.URL.Path
		strippedPath:=strings.TrimPrefix(originalPath,pathPrefix)

		r.URL.Host=target.Host
		r.URL.Path=target.Path+strippedPath
		
		r.Host=target.Host

		if userId,ok:=r.Context().Value("userID").(string); ok{
			r.Header.Set("X-User-ID",userId)
		}
	}

	return proxy.ServeHTTP
}