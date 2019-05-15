package xcos

import (
	"context"
	"io"
	"net/http"
	"strings"

	cos "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/webee/x/xfile"
)

// CosFileResource Cos文件资源
type CosFileResource struct {
	client *Client
}

const (
	// CosPathPrefix cos路径前缀
	CosPathPrefix = "cos:/"
)

// NewCosFileResourceHandler 新建cos文件处理器
func NewCosFileResourceHandler(client *Client) *CosFileResource {
	return &CosFileResource{client}
}

// CanHandle 文件资源处理接口
func (h *CosFileResource) CanHandle(path string) bool {
	return strings.HasPrefix(path, CosPathPrefix)
}

func (h *CosFileResource) path(path string) string {
	return strings.TrimLeft(path, CosPathPrefix)
}

// OpenForRead 文件资源处理接口
func (h *CosFileResource) OpenForRead(path string) (reader io.ReadCloser, err error) {
	var (
		ctx  = context.Background()
		resp *cos.Response
	)

	path = h.path(path)
	if resp, err = h.client.COS().Object.Get(ctx, path, nil); err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			err = xfile.ErrFileNotExits
		}
		return
	}
	reader = resp.Body
	return
}

// CreateFile 文件资源处理接口
func (h *CosFileResource) CreateFile(path string, reader io.Reader) (etag string, err error) {
	path = h.path(path)
	return h.client.Create(path, reader)
}

// ContentLength 文件资源处理接口
func (h *CosFileResource) ContentLength(path string) (length int64, err error) {
	var (
		meta *ObjectMetaData
	)

	path = h.path(path)

	if meta, err = h.client.Head(path); err != nil {
		return
	}

	length = meta.ContentLength
	return
}

// EntityTag 文件资源处理接口
func (h *CosFileResource) EntityTag(path string) (etag string, err error) {
	var (
		meta *ObjectMetaData
	)

	path = h.path(path)
	if meta, err = h.client.Head(path); err != nil {
		return
	}

	etag = meta.ETag
	return
}

// IsModified 文件资源处理接口
func (h *CosFileResource) IsModified(path string, etag string) (ok bool, err error) {
	return xfile.FileIsModified(h, path, etag)
}
