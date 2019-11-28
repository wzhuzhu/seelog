package seelog

import (
	"errors"
	"fmt"
)

type slog struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

var slogs = []slog{}

// SeeAdd 添加需要监视的文件
func SeeAdd(name, path string) (err error) {

	if name == "" || path == "" {
		err = errors.New("log名称或者路径不可为空")
		return
	}

	for _, sl := range slogs {
		if sl.Name == name {
			err = fmt.Errorf("log名称 %s 已存在,不可重复", name)
			return
		}
	}
	slogs = append(slogs, slog{name, path})
	return
}

// Serve 启动服务
func Serve(port int) (err error) {

	if port < 0 || port > 65535 {
		err = errors.New("端口号不符合规范，port(0,65535)")
		return
	}

	if len(slogs) < 1 {
		err = errors.New("至少监听一个日志文件,请使用 seelog.SeeAdd(name,path string)")
		return
	}
	// 开启socket管理器
	go manager.start()

	// 监控文件
	go monitor()

	// 开启httpServer
	go server(port)

	return
}
