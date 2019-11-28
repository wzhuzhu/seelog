package seelog

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	auth "github.com/abbot/go-http-auth"
	"golang.org/x/net/websocket"
)

// 开启 httpServer
func server(port int) (err error) {

	defer func() {
		if err := recover(); err != nil {
			printError(errors.New("server panic"))
		}
	}()

	// socket链接
	http.Handle("/ws", websocket.Handler(genConn))

	_authlog := auth.NewBasicAuthenticator("frnd log", secret)
	http.HandleFunc("/log", auth.JustCheck(_authlog, logPage))

	_authhitory := auth.NewBasicAuthenticator("frnd log history", secret)
	if logDir != "" {
		http.HandleFunc("/history/", auth.JustCheck(_authhitory, handleFileServer(logDir, "/history/")))
	}

	// 访问页面
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		showPage(writer, page403, nil)
	})
	if tlsCrt != "" && tlsKey != "" {
		err = http.ListenAndServeTLS(fmt.Sprintf(":%d", port), tlsCrt, tlsKey, nil)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}
	return
}

func logPage(w http.ResponseWriter, r *http.Request) {
	showPage(w, pageIndex, slogs)
}

// 输出page
func showPage(writer http.ResponseWriter, page string, data interface{}) {
	t, err := template.ParseFiles(page)
	if err != nil {
		printError(err)
	}
	t.Execute(writer, data)
}

// 创建client对象
func genConn(ws *websocket.Conn) {
	client := &client{time.Now().String(), ws, make(chan msg, 1024), slogs[0].Name}
	manager.register <- client
	go client.read()
	client.write()
}
