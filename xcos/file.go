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

// OpenForRead 文件资源处理接口
func (h *CosFileResource) OpenForRead(path string) (reader io.ReadCloser, err error) {
	var (
		ctx  = context.Background()
		resp *cos.Response
	)

	path = strings.TrimLeft(path, CosPathPrefix)
	if resp, err = h.client.COS().Object.Get(ctx, path, nil); err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			err = xfile.ErrFileNotExits
		}
		return
	}
	reader = resp.Body
	return
}

// EntityTag 文件资源处理接口
func (h *CosFileResource) EntityTag(path string) (etag string, err error) {
	var (
		ctx  = context.Background()
		resp *cos.Response
	)

	path = strings.TrimLeft(path, CosPathPrefix)
	if resp, err = h.client.COS().Object.Head(ctx, path, nil); err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			err = xfile.ErrFileNotExits
		}
		return
	}

	etag = resp.Header.Get("ETag")
	return
}

// IsModified 文件资源处理接口
func (h *CosFileResource) IsModified(path string, etag string) (ok bool, err error) {
	return xfile.FileIsModified(h, path, etag)
}
