package xfile

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
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

// CreateFile 文件资源处理接口
func (h *HTTPFileResourceHandler) CreateFile(path string, reader io.Reader) (etag string, err error) {
	err = ErrNotImplemented
	return
}

// ContentLength 文件资源处理接口
func (h *HTTPFileResourceHandler) ContentLength(path string) (length int64, err error) {
	var (
		resp *http.Response
	)
	if resp, err = http.Head(path); err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			err = ErrFileNotExits
		}
		return
	}

	lengthStr := resp.Header.Get("Content-Length")
	if lengthStr == "" {
		err = ErrNotImplemented
		return
	}
	return strconv.ParseInt(lengthStr, 10, 64)
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
