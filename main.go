package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	urlFix       = regexp.MustCompile(`^(https?:/)`)
	urlMatcher01 = regexp.MustCompile(`^(?:(?:https?://)?github\.com/)?[^/]+?/[^/]+?/(?:releases|archive)/.*$`)
	urlMatcher02 = regexp.MustCompile(`^(?:(?:https?://)?github\.com/)?[^/]+?/[^/]+?/(?:blob|raw)/.*$`)
	urlMatcher03 = regexp.MustCompile(`^(?:(?:https?://)?github\.com/)?[^/]+?/[^/]+?/(?:info|git-).*$`)
	urlMatcher04 = regexp.MustCompile(`^(?:(?:https?://)?github\.com/)?[^/]+?/[^/]+?/tags.*$`)
	urlMatcher10 = regexp.MustCompile(`^(?:https?://)?raw\.(?:githubusercontent|github)\.com/(?P<author>.+?)/(?P<repo>.+?)/.+?/.+$`)
	urlMatcher20 = regexp.MustCompile(`^(?:https?://)?gist\.(?:githubusercontent|github)\.com/(?P<author>.+?)/.+?/.+$`)
)

func matchUrl(u string) bool {
	for _, exp := range []*regexp.Regexp{urlMatcher01, urlMatcher02, urlMatcher03, urlMatcher04, urlMatcher10, urlMatcher20} {
		if exp.MatchString(u) {
			return true
		}
	}
	return false
}

func handler(w http.ResponseWriter, req *http.Request) {
	u := req.URL.RequestURI()[1:]

	// 修正 URL，由于URL在进行处理时，会将两个url合并为一个，确保路径开头有两个斜杠
	if strings.HasPrefix(u, "http") {
		u = urlFix.ReplaceAllString(u, "$1/")
	}

	// preflight
	if preflight(w, req) {
		return
	}

	// 不满足要求的URL不会进行处理
	if !matchUrl(u) {
		w.WriteHeader(http.StatusForbidden)
	}

	// 补充完整URL信息
	if urlMatcher01.MatchString(u) || urlMatcher02.MatchString(u) || urlMatcher03.MatchString(u) || urlMatcher04.MatchString(u) {
		if !strings.HasPrefix(u, "http") {
			if strings.HasPrefix(u, "github.com") {
				u = "https://" + u
			} else {
				u = "https://github.com/" + u
			}
		}
	} else if !strings.HasPrefix(u, "http") {
		u = "https://" + u
	}

	fmt.Println("Received URL:", u)

	// 检查Url是否有效
	proxyUrl, err := url.Parse(u)
	if err != nil {
		http.Error(w, "Failed to parse url.", http.StatusInternalServerError)
		return
	}

	// 转发请求, 获取文件
	proxy(w, &http.Request{
		Method: req.Method,
		URL:    proxyUrl,
		Header: req.Header,
		Body:   req.Body,
	})
}

func preflight(w http.ResponseWriter, req *http.Request) bool {
	if req.Method == "OPTIONS" && req.Header.Get("Access-Control-Request-Headers") != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Max-Age", "600")
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func proxy(w http.ResponseWriter, req *http.Request) {
	httpClient := http.Client{}
	// 执行请求
	res, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, "failed to fetch from upstream", http.StatusBadGateway)
		return
	}
	resHeader := res.Header

	// 设置响应头部
	resHeader.Set("Access-Control-Expose-Headers", "*")
	resHeader.Set("Access-Control-Allow-Origin", "*")

	resHeader.Del("Content-Security-Policy")
	resHeader.Del("Content-Security-Policy-Report-Only")
	resHeader.Del("Clear-Site-Data")

	// 写入响应
	for key, values := range resHeader {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(res.StatusCode)
	_, _ = io.Copy(w, res.Body)
}

func main() {
	http.HandleFunc("/", handler)
	_ = http.ListenAndServe(":8080", nil)
}
