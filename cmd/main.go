package main

import (
	"log"
	"net/http"

	"golang.org/x/net/webdav"
)

func main() {
	// 设置要共享的目录
	dir := "/path/to/your/directory"

	// 创建一个WebDAV文件服务器
	fs := &webdav.Handler{
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
	}

	// 启动WebDAV服务器
	port := ":8080"
	serverAddr := "0.0.0.0" + port
	log.Printf("WebDAV server is listening on %s...\n", serverAddr)
	err := http.ListenAndServe(serverAddr, fs)
	if err != nil {
		log.Fatal("Error:", err)
	}
}
