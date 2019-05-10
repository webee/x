package xfile

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HTTPFileResourceHandler Http文件资源
type HTTPFileResourceHandler struct{}

// CanHandle 文件资源处理接口
func (h *HTTPFileResourceHandler) CanHandle(path string) bool {
	var (
		u   *url.URL
		err error
	)

	if u, err = url.Parse(path); err != nil {
		return false
	}

	return strings.HasPrefix(u.Scheme, "http")
}

// OpenForRead 文件资源处理接口
func (h *HTTPFileResourceHandler) OpenForRead(path string) (reader io.ReadCloser, err error) {
	var (
		resp *http.Response
	)

	if resp, err = http.Get(path); err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			err = ErrFileNotExits
		}
		return
	}
	reader = resp.Body
	return
}

// EntityTag 文件资源处理接口
func (h *HTTPFileResourceHandler) EntityTag(path string) (etag string, err error) {
	var (
		resp *http.Response
	)
	if resp, err = http.Head(path); err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			err = ErrFileNotExits
		}
		return
	}

	etag = resp.Header.Get("ETag")
	if etag == "" {
		etag = resp.Header.Get("Last-Modified")
	}
	return
}

// IsModified 文件资源处理接口
func (h *HTTPFileResourceHandler) IsModified(path string, etag string) (ok bool, err error) {
	return FileIsModified(h, path, etag)
}
