package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coconutLatte/volume-adaptor/webdav"
)

func main() {
	var err error
	adaptor, err := webdav.NewYSStorageAdaptor()
	if err != nil {
		panic(fmt.Sprintf("new ys storage adaptor failed, %v", err))
	}

	// 创建一个WebDAV文件服务器
	fs := &webdav.Handler{
		FileSystem: adaptor,
		//FileSystem: webdav.Dir("/"),
		LockSystem: webdav.NewMemLS(),
	}

	// 启动WebDAV服务器
	port := ":8080"
	serverAddr := "0.0.0.0" + port
	log.Printf("WebDAV server is listening on %s...\n", serverAddr)
	err = http.ListenAndServe(serverAddr, fs)
	if err != nil {
		log.Fatal("Error:", err)
	}
}
