package webdav

import (
	"context"
	"errors"
	"fmt"
	"io"
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

func (s *YSStorageAdaptor) Mkdir(_ context.Context, name string, _ os.FileMode) error {
	logging.Default().Infof("mkdir, name: %s", name)

	return s.cli.Mkdir(name)
}

type ysStorageFile struct {
	cli *openapi.Client

	path      string
	fileIndex int64
}

func newYsStorageFile(cli *openapi.Client, path string) (*ysStorageFile, error) {
	_, err := cli.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if err = cli.Create(path); err != nil {
			return nil, os.ErrInvalid
		}
	}

	f := &ysStorageFile{
		cli:  cli,
		path: path,
	}

	return f, nil
}

func (f *ysStorageFile) Read(p []byte) (n int, err error) {
	logging.Default().Info("ysStorageFile read")

	// TODO make it cache able
	fi, err := f.cli.Stat(f.path)
	if err != nil {
		return 0, os.ErrInvalid
	}
	if fi.IsDir() {
		return 0, os.ErrInvalid
	}

	data, err := f.cli.ReadAt(f.path, f.fileIndex, int64(len(p)))
	if err != nil {
		logging.Default().Errorf("read file failed, %v", err)
		return -1, err
	}
	copy(p, data)

	f.fileIndex += int64(len(data))

	logging.Default().Infof("read result: %s", string(p))

	return len(data), nil
}

func (f *ysStorageFile) Seek(offset int64, whence int) (int64, error) {
	logging.Default().Infof("ysStorageFile seek, offset: %d, whence: %d", offset, whence)

	// TODO make it cache able
	fi, err := f.cli.Stat(f.path)
	if err != nil {
		return 0, os.ErrInvalid
	}

	npos := f.fileIndex
	switch whence {
	case io.SeekStart:
		npos = offset
	case io.SeekCurrent:
		npos += offset
	case io.SeekEnd:
		npos = fi.Size() + offset
	default:
		npos = -1
	}

	if npos < 0 {
		return 0, os.ErrInvalid
	}
	f.fileIndex = npos

	return f.fileIndex, nil
}

func (f *ysStorageFile) Readdir(count int) ([]fs.FileInfo, error) {
	logging.Default().Info("ysStorageFile readdir")

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
	logging.Default().Info("ysStorageFile stat")

	fi, err := f.cli.Stat(f.path)
	if err != nil {
		logging.Default().Error(err)
		return nil, err
	}

	return fi, nil
}

func (f *ysStorageFile) Write(p []byte) (n int, err error) {
	logging.Default().Info("ysStorageFile write")

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

func (s *YSStorageAdaptor) OpenFile(_ context.Context, name string, _ int, _ os.FileMode) (File, error) {
	logging.Default().Infof("openFile, name: %s", name)

	f, err := newYsStorageFile(s.cli, name)
	if err != nil {
		return nil, os.ErrInvalid
	}

	return f, nil
}

func (s *YSStorageAdaptor) RemoveAll(_ context.Context, name string) error {
	logging.Default().Infof("removeAll, name: %s", name)

	return s.cli.RemoveAll(name)
}

func (s *YSStorageAdaptor) Rename(_ context.Context, oldName, newName string) error {
	logging.Default().Infof("rename, oldName: %s, newName: %s", oldName, newName)

	return s.cli.Rename(oldName, newName)
}

func (s *YSStorageAdaptor) Stat(_ context.Context, name string) (os.FileInfo, error) {
	logging.Default().Infof("stat, name: %s", name)

	fi, err := s.cli.Stat(name)
	if err != nil {
		logging.Default().Errorf("call stat api failed, %v", err)
		return nil, err
	}

	return fi, nil
}
