package xfile

import (
	"errors"
	"io"
)

// FileResourceHandler 文件资源处理器接口
type FileResourceHandler interface {
	CanHandle(path string) bool                        // 是否可以处理此路径
	OpenForRead(path string) (io.ReadCloser, error)    // 打开文件，得到reader
	EntityTag(path string) (string, error)             // 获取文件唯一标签
	IsModified(path string, etag string) (bool, error) // 是否修改了
}

// errors
var (
	ErrFileNotExits     = errors.New("file not exists")
	ErrCannotHandleFile = errors.New("can not handle file")
)

const (
	// EmptyETag 空的etag
	EmptyETag = ""
)

// FileIsModified 文件是否修改了
func FileIsModified(handler FileResourceHandler, path string, etag string) (ok bool, err error) {
	var (
		curETag string
	)
	if curETag, err = handler.EntityTag(path); err != nil {
		return
	}
	if curETag == EmptyETag {
		// 没有有效的etag, 则认为修改了
		ok = true
		return
	}
	ok = curETag != etag
	return
}

// MultiFileResourceHandler 多重文件资源处理器
type MultiFileResourceHandler struct {
	handlers []FileResourceHandler
}

// NewMultiFileResourceHandler 新建多重文件资源处理器
func NewMultiFileResourceHandler(handlers ...FileResourceHandler) (res *MultiFileResourceHandler) {
	res = new(MultiFileResourceHandler)
	res.handlers = append(res.handlers, handlers...)
	return
}

// RegisterHandler 注册处理器
func (h *MultiFileResourceHandler) RegisterHandler(handler FileResourceHandler) {
	handlers := []FileResourceHandler{handler}
	handlers = append(handlers, h.handlers...)
	h.handlers = handlers
}

// CanHandle 文件资源处理接口
func (h *MultiFileResourceHandler) CanHandle(path string) bool {
	_, err := h.canHandle(path)
	return err != nil
}

func (h *MultiFileResourceHandler) canHandle(path string) (handler FileResourceHandler, err error) {
	for _, handler = range h.handlers {
		if handler.CanHandle(path) {
			return
		}
	}
	err = ErrCannotHandleFile
	return
}

// OpenForRead 文件资源处理接口
func (h *MultiFileResourceHandler) OpenForRead(path string) (reader io.ReadCloser, err error) {
	var (
		handler FileResourceHandler
	)
	if handler, err = h.canHandle(path); err != nil {
		return
	}

	return handler.OpenForRead(path)
}

// EntityTag 文件资源处理接口
func (h *MultiFileResourceHandler) EntityTag(path string) (etag string, err error) {
	var (
		handler FileResourceHandler
	)
	if handler, err = h.canHandle(path); err != nil {
		return
	}

	return handler.EntityTag(path)
}

// IsModified 文件资源处理接口
func (h *MultiFileResourceHandler) IsModified(path string, etag string) (ok bool, err error) {
	return FileIsModified(h, path, etag)
}