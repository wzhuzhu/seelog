package seelog

import (
	"errors"
	"log"
)

type slog struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

var slogs = []slog{}


// 启动seelog
func See(name,path string) {

	if name == "" || path == "" {
		log.Println(errors.New("名称或路径不能为空"))
		return
	}
	slogs = append(slogs, slog{name,path})
}

// 开始监控
func Serve(port int,password string)  {

	if port < 0 || port > 65535 {
		log.Println(errors.New("端口号不正确"))
		return
	}
	// 开启socket管理器
	go manager.start()

	// 监控文件
	go monitor()

	// 开启httpServer
	go server(port,password)
}
