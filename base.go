package seelog

import (
	"fmt"
	"log"
	"net/http"
	"os"

	auth "github.com/abbot/go-http-auth"
)

var authUser, authPass, authMD5, tlsCrt, tlsKey, pageDir, logDir, pageIndex, page403 string

func init() {
	authUser = "frnd"
	authPass = "frnd21*71Log"
	authMD5 = string(auth.MD5Crypt([]byte(authPass), []byte("xQiYxTAr"), []byte("$9$")))
	pageDir = "/tmp"
}

// Config 设置配置参数
func Config(user, pass, tlscrt, tlskey, pagedir, logdir string) (err error) {
	authUser = user
	authPass = pass
	authMD5 = string(auth.MD5Crypt([]byte(pass), []byte("xQiYxTAr"), []byte("$9$")))
	tlsCrt = tlscrt
	tlsKey = tlskey
	pageDir = pagedir
	logDir = logdir
	pageIndex = pageDir + "/index.html"
	page403 = pageDir + "/403.html"

	if checkFileExist(tlsCrt) == false {
		err = fmt.Errorf("tlsCrt: %s, file not exist", tlsCrt)
		return
	}
	if checkFileExist(tlsKey) == false {
		err = fmt.Errorf("tlsKey: %s, file not exist", tlsKey)
		return
	}
	if checkFileExist(pageIndex) == false {
		err = fmt.Errorf("pageIndex: %s, file not exist", pageIndex)
		return
	}
	if checkFileExist(page403) == false {
		err = fmt.Errorf("page403: %s, file not exist", page403)
		return
	}
	return
}

// secret 验证密码
func secret(u, realm string) string {
	if u != authUser {
		return ""
	}
	return authMD5
}

func handleFileServer(dir, prefix string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))
	realHandler := http.StripPrefix(prefix, fs).ServeHTTP
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.URL)
		realHandler(w, req)
	}
}

// checkFileExist 检查文件是否存在
func checkFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return true
}
