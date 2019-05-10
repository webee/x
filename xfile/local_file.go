package xfile

import (
	"io"
	"net/url"
	"os"
)

// LocalFileResourceHandler 本地文件资源
type LocalFileResourceHandler struct{}

// CanHandle 文件资源处理接口
func (h *LocalFileResourceHandler) CanHandle(path string) bool {
	var (
		u   *url.URL
		err error
	)

	if u, err = url.Parse(path); err != nil {
		return false
	}

	return u.Scheme == "" && u.Host == ""
}

// OpenForRead 文件资源处理接口
func (h *LocalFileResourceHandler) OpenForRead(path string) (reader io.ReadCloser, err error) {
	if reader, err = os.Open(path); err != nil {
		if os.IsNotExist(err) {
			err = ErrFileNotExits
		}
		return
	}
	return
}

// EntityTag 文件资源处理接口
func (h *LocalFileResourceHandler) EntityTag(path string) (etag string, err error) {
	var (
		stat os.FileInfo
	)
	if stat, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = ErrFileNotExits
		}
		return
	}

	etag = stat.ModTime().String()
	return
}

// IsModified 文件资源处理接口
func (h *LocalFileResourceHandler) IsModified(path string, etag string) (ok bool, err error) {
	return FileIsModified(h, path, etag)
}
