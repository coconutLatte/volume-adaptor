package webdav

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"git.yuansuan.cn/go-kit/logging"

	"github.com/coconutLatte/volume-adaptor/openapi"
)

type YSStorageAdaptor struct {
	cli *openapi.Client
}

func NewYSStorageAdaptor() (*YSStorageAdaptor, error) {
	cli, err := openapi.NewClient("", "", "")
	if err != nil {
		return nil, fmt.Errorf("new openapi client failed, %w", err)
	}

	return &YSStorageAdaptor{
		cli: cli,
	}, nil
}

func (s *YSStorageAdaptor) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	logging.Default().Info("mkdir")

	return s.cli.Mkdir(name)
}

type ysStorageFile struct {
	cli *openapi.Client
	f   *os.File

	path      string
	fileIndex int64
}

func (f *ysStorageFile) Read(p []byte) (n int, err error) {
	data, err := f.cli.ReadAt(f.path, f.fileIndex, int64(len(p)))
	if err != nil {
		logging.Default().Errorf("read file failed, %v", err)
		return -1, err
	}
	p = data
	f.fileIndex += int64(len(data))

	return len(data), nil
}

func (f *ysStorageFile) Seek(offset int64, whence int) (int64, error) {
	// TODO
	return f.f.Seek(offset, whence)
}

func (f *ysStorageFile) Readdir(count int) ([]fs.FileInfo, error) {
	// call stat
	fi, err := f.cli.Stat(f.path)
	if err != nil {
		logging.Default().Errorf("call stat failed, %v", err)
		return nil, err
	}

	if !fi.IsDir() {
		return nil, os.ErrInvalid
	}

	fis, err := f.cli.LsWithPage(f.path, int64(count))
	if err != nil {
		logging.Default().Error(err)
		return nil, err
	}

	return fis, nil
}

func (f *ysStorageFile) Stat() (fs.FileInfo, error) {
	fi, err := f.cli.Stat(f.path)
	if err != nil {
		logging.Default().Error(err)
		return nil, err
	}

	return fi, nil
}

func (f *ysStorageFile) Write(p []byte) (n int, err error) {
	lenP := len(p)

	err = f.cli.WriteAt(f.path, p, f.fileIndex)
	if err != nil {
		logging.Default().Error(err)
		return -1, err
	}
	f.fileIndex += int64(lenP)

	return lenP, nil
}

func (f *ysStorageFile) Close() error {
	// do nothing
	return nil
}

func (s *YSStorageAdaptor) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (File, error) {
	return &ysStorageFile{
		cli:  s.cli,
		path: name,
	}, nil
}

func (s *YSStorageAdaptor) RemoveAll(ctx context.Context, name string) error {
	return nil
}

func (s *YSStorageAdaptor) Rename(ctx context.Context, oldName, newName string) error {
	return nil
}

func (s *YSStorageAdaptor) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	fi, err := s.cli.Stat(name)
	if err != nil {
		logging.Default().Errorf("call stat api failed, %v", err)
		return nil, err
	}

	return fi, nil
}
