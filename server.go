package seelog

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	PageIndex = "index.html"
	Page403   = "403.html"
)

// 开启 httpServer
func server(port int, password string) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("[seelog] error:%+v\n", err)
		}
	}()

	// socket链接
	http.Handle("/ws", websocket.Handler(genConn))

	// 访问页面
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if !(strings.Replace(request.RequestURI, "/", "", -1) == password) {
			showPage(writer, Page403, nil)
			return
		}
		showPage(writer, PageIndex, slogs)
	})
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// 输出page
func showPage(writer http.ResponseWriter, page string, data interface{}) {
	_, currentfile, _, _ := runtime.Caller(0) // 忽略错误
	filename := path.Join(path.Dir(currentfile), page)
	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Println(err)
	}
	t.Execute(writer, data)
}

// 创建client对象
func genConn(ws *websocket.Conn) {
	var defaultLog = ""
	if len(slogs) > 0 {
		defaultLog = slogs[0].Name
	}
	client := &client{time.Now().String(), ws, make(chan msg, 1024), defaultLog}
	manager.register <- client
	go client.recv()
	client.write()
}
