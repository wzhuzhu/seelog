package seelog

import (
	"github.com/hpcloud/tail"
	"log"
)

type msg struct {
	LogName string `json:"logName"`
	Data	string `json:"data"`
}

// 监控日志文件
func monitor() {

	for _,sl := range slogs {
		go func(sl slog) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("[seelog] error:%+v", err)
				}
			}()

			log.Println("开始进行日志监控",sl.Name, sl.Path)

			t, _ := tail.TailFile(sl.Path, tail.Config{Follow: true})
			for line := range t.Lines {
				manager.broadcast <- msg{sl.Name,line.Text}
			}
		}(sl)
	}

}
