package openapi

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"git.yuansuan.cn/go-kit/logging"
	openys "git.yuansuan.cn/openapi-go"
	"git.yuansuan.cn/openapi-go/credential"
	v20230530 "git.yuansuan.cn/project-root-api/schema/v20230530"
)

const (
	storageEndpoint = "http://10.0.4.191:8899"
	defaultYsId     = "4YZM5xWH91Y"
)

type Client struct {
	ysId string
	base *openys.Client
}

func NewClient(ysId, accessKey, accessSecret string) (*Client, error) {
	if ysId == "" {
		ysId = defaultYsId
	}

	cli, err := openys.NewClient(
		credential.NewCredential(accessKey, accessSecret),
		openys.WithBaseURL(storageEndpoint),
		openys.WithRetryTimes(1),
		openys.WithRetryInterval(time.Millisecond),
	)
	if err != nil {
		return nil, fmt.Errorf("create openapi client failed")
	}

	return &Client{
		ysId: ysId,
		base: cli,
	}, nil
}

func (c *Client) Mkdir(dir string) error {
	_, err := c.base.Storage.Mkdir(
		c.base.Storage.Mkdir.Path(c.filePath(dir)),
		c.base.Storage.Mkdir.IgnoreExist(true),
	)
	return err
}

func (c *Client) ReadAt(path string, offset, length int64) ([]byte, error) {
	resp, err := c.base.Storage.ReadAt(
		c.base.Storage.ReadAt.Path(c.filePath(path)),
		c.base.Storage.ReadAt.Offset(offset),
		c.base.Storage.ReadAt.Length(length),
	)
	if err != nil {
		logging.Default().Errorf("call read at api failed, %v", err)
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) WriteAt(path string, p []byte, offset int64) error {
	_, err := c.base.Storage.WriteAt(
		c.base.Storage.WriteAt.Path(c.filePath(path)),
		c.base.Storage.WriteAt.Offset(offset),
		c.base.Storage.WriteAt.Length(int64(len(p))),
		c.base.Storage.WriteAt.Data(bytes.NewBuffer(p)),
	)
	if err != nil {
		logging.Default().Errorf("call write at api failed, %v", err)
		return err
	}

	return nil
}

type fileInfo struct {
	openapiFi *v20230530.FileInfo
}

func (fi *fileInfo) Name() string {
	return fi.openapiFi.Name
}

func (fi *fileInfo) Size() int64 {
	return fi.openapiFi.Size
}

func (fi *fileInfo) Mode() os.FileMode {
	return os.FileMode(fi.openapiFi.Mode)
}

func (fi *fileInfo) ModTime() time.Time {
	return fi.openapiFi.ModTime
}

func (fi *fileInfo) IsDir() bool {
	return fi.openapiFi.IsDir
}

func (fi *fileInfo) Sys() any {
	return nil
}

func (c *Client) Stat(path string) (os.FileInfo, error) {
	logging.Default().Info("stat")
	logging.Default().Info("path: ", path)

	resp, err := c.base.Storage.Stat(
		c.base.Storage.Stat.Path(c.filePath(path)),
	)
	if err != nil {
		return nil, os.ErrNotExist

		//errResp, ok := err.(*common.ErrorResp)
		//if ok {
		//	if err = jsoniter.Unmarshal(errResp.RawBody(), resp); err != nil {
		//		return nil, err
		//	}
		//
		//	logging.Default().Errorf("%#v", resp)
		//	if resp.Response.ErrorCode == "PathNotFound" {
		//		//return nil, os.ErrNotExist
		//		return nil, &fs.PathError{}
		//	}
		//}
		//
		//return nil, err
	}
	if resp == nil || resp.Data == nil || resp.Data.File == nil {
		return nil, fmt.Errorf("invalid stat response data")
	}

	return &fileInfo{openapiFi: resp.Data.File}, nil
}

func (c *Client) filePath(path string) string {
	res := "/" + c.ysId + path

	return res
}

func (c *Client) LsWithPage(path string, count int64) ([]os.FileInfo, error) {
	if count == 0 {
		count = 100
	}

	resp, err := c.base.Storage.LsWithPage(
		c.base.Storage.LsWithPage.Path(c.filePath(path)),
		c.base.Storage.LsWithPage.PageSize(count),
	)
	if err != nil {
		logging.Default().Error(err)
		return nil, err
	}
	if resp == nil || resp.Data == nil {
		logging.Default().Errorf("invalid resp")
		return nil, fmt.Errorf("invalid resp")
	}

	fis := make([]os.FileInfo, 0)
	for _, fi := range resp.Data.Files {
		fis = append(fis, &fileInfo{
			openapiFi: fi,
		})
	}

	return fis, nil
}

func (c *Client) Rename(src, dst string) error {
	logging.Default().Infof("rename, src [%s], dst [%s]", src, dst)

	_, err := c.base.Storage.Mv(
		c.base.Storage.Mv.Src(c.filePath(src)),
		c.base.Storage.Mv.Dest(c.filePath(dst)),
	)
	if err != nil {
		return os.ErrInvalid
	}

	return nil
}

func (c *Client) RemoveAll(path string) error {
	logging.Default().Infof("removeAll, path [%s]", path)

	_, err := c.base.Storage.Rm(
		c.base.Storage.Rm.Path(c.filePath(path)),
		c.base.Storage.Rm.IgnoreNotExist(true),
	)
	if err != nil {
		return os.ErrInvalid
	}

	return nil
}
