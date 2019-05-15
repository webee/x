package xcos

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/hyacinthus/x/object"
	gocos "github.com/mozillazg/go-cos"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/webee/x/xfile"
)

// ObjectMetaData 对象元数据
type ObjectMetaData struct {
	ContentType   string
	ContentLength int64
	ETag          string
	LastModified  time.Time
}

// Client cos client
type Client struct {
	*object.Client
	config *object.Config
	cos    *gocos.Client
}

// New 新建Client
func New(client *object.Client, config *object.Config) *Client {
	baseURL, _ := gocos.NewBaseURL(gocos.NewBucketURL(config.Bucket, config.AppID, config.Region, true).String())
	return &Client{
		Client: client,
		config: config,
		cos: gocos.NewClient(baseURL, &http.Client{
			Transport: &gocos.AuthorizationTransport{
				SecretID:  config.SecretID,
				SecretKey: config.SecretKey,
				Expire:    time.Hour,
			},
		}),
	}
}

// GetObjectPresignedURL 得到获取对象的预签名链接
func (c *Client) GetObjectPresignedURL(key string, expire time.Duration) (*url.URL, error) {
	var (
		ctx = context.Background()
	)

	auth := gocos.Auth{
		SecretID:  c.config.SecretID,
		SecretKey: c.config.SecretKey,
		Expire:    expire,
	}
	return c.cos.Object.PresignedURL(ctx, http.MethodGet, key, auth, nil)
}

// Create 创建文件
func (c *Client) Create(key string, r io.Reader) (etag string, err error) {
	var (
		ctx  = context.Background()
		resp *cos.Response
	)

	if resp, err = c.COS().Object.Put(ctx, key, r, &cos.ObjectPutOptions{}); err != nil {
		return
	}

	etag = resp.Header.Get("Etag")
	return
}

// Head 获取 cos 对象元数据
func (c *Client) Head(key string) (*ObjectMetaData, error) {
	var (
		err error
		ctx = context.Background()
	)

	resp, err := c.COS().Object.Head(ctx, key, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			err = xfile.ErrFileNotExits
		}
		return nil, err
	}

	meta := new(ObjectMetaData)
	meta.ContentType = resp.Header.Get("Content-Type")
	if len, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64); err == nil {
		meta.ContentLength = len
	}
	meta.ETag = resp.Header.Get("ETag")
	if t, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified")); err == nil {
		meta.LastModified = t
	}
	return meta, nil
}

// Exists 获取 cos 对象元数据
func (c *Client) Exists(key string) bool {
	var (
		err error
		ctx = context.Background()
	)

	_, err = c.COS().Object.Head(ctx, key, nil)
	return err == nil
}

// DeleteMulti 删除多个对象
func (c *Client) DeleteMulti(keys ...string) (*cos.ObjectDeleteMultiResult, error) {
	var (
		err error
		ctx = context.Background()
	)

	objs := make([]cos.Object, len(keys))
	for i := range objs {
		objs[i] = cos.Object{Key: keys[i]}
	}

	result, _, err := c.COS().Object.DeleteMulti(ctx, &cos.ObjectDeleteMultiOptions{Objects: objs})
	return result, err
}

// BucketObjectList 存储桶对象列表
type BucketObjectList struct {
	client *Client
	prefix string
}

// NewBucketObjectList 新建一个桶对象列表
func (c *Client) NewBucketObjectList(prefix string) *BucketObjectList {
	return &BucketObjectList{
		client: c,
		prefix: prefix,
	}
}

// Iterator 返回对象迭代器
func (l *BucketObjectList) Iterator() (<-chan *cos.Object, func(), func() error) {
	var (
		err  error
		stop = make(chan struct{})
		data = make(chan *cos.Object)
		once = sync.Once{}
	)

	stopFunc := func() {
		once.Do(func() {
			close(stop)
		})
	}

	errFunc := func() error {
		return err
	}

	go func() {
		var (
			result *cos.BucketGetResult
			ctx    = context.Background()
		)
	LOOP:
		for result == nil || result.IsTruncated {
			select {
			case <-stop:
				goto DONE
			default:
				opt := &cos.BucketGetOptions{
					Prefix: l.prefix,
				}
				if result != nil {
					opt.Marker = result.NextMarker
				}
				result, _, err = l.client.COS().Bucket.Get(ctx, opt)
				if err != nil {
					break LOOP
				}
				for i := range result.Contents {
					select {
					case <-stop:
						goto DONE
					default:
						data <- &result.Contents[i]
					}
				}
			}
		}
	DONE:
		close(data)
		stopFunc()
	}()

	return data, stopFunc, errFunc
}
