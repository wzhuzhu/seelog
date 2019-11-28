package seelog_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/wzhuzhu/seelog"
)

const (
	DebugLog = "debug.log"
	ErrLog   = "err.log"
)

func TestSee(t *testing.T) {

	// 测试
	seelog.See("错误日志", ErrLog)
	seelog.See("调试日志", DebugLog)
	seelog.Serve(9000)

	// 模拟服务输出日志
	go printLog("调试日志", DebugLog)
	go printLog("错误日志", ErrLog)
	select {}
}

func printLog(name, path string) {

	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return
	}

	for t := range time.Tick(time.Second * 1) {
		testLog := fmt.Sprintf("「%s」[%s]\n", name, t.String())
		_, err := f.WriteString(testLog)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
