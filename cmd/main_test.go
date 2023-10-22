package main

import (
	"fmt"
	"os"
	"testing"

	openys "git.yuansuan.cn/openapi-go"
	"git.yuansuan.cn/openapi-go/credential"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	f, err := os.OpenFile("../test_open_file", os.O_RDWR, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf1 := make([]byte, 1)
	n, err := f.Read(buf1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("1: ", string(buf1[:n]))
	fmt.Println(n)

	n, err = f.Write([]byte("zzz"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(n)

	n, err = f.Read(buf1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("2: ", string(buf1[:n]))
	fmt.Println(n)
}

const ysid = "4YZM5xWH91Y"

func TestReadApi(t *testing.T) {
	cli, err := openys.NewClient(
		credential.NewCredential("", ""),
		openys.WithBaseURL("http://10.0.4.191:8899"),
	)
	assert.NoError(t, err)

	resp, err := cli.Storage.ReadAt(
		cli.Storage.ReadAt.Path("/"+ysid+"/abc/test.txt"),
		cli.Storage.ReadAt.Length(1024),
	)
	assert.NoError(t, err)

	fmt.Println(string(resp.Data))
	fmt.Println(len(resp.Data))
}
