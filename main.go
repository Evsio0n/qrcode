package main

import (
	"./module/library/log"
	"./module/library/qrcode"
	"context"
	"fmt"
	"github.com/medivh-jay/daemon"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func qr(w http.ResponseWriter, r *http.Request) {
	var c string
	var wi int
	var h int
	var hs string
	var ws string
	var key string
	var count int
	var problem int
	var oc int
	var err error
	_ = r.ParseForm() //解析参数
	count = 0
	for k, v := range r.Form {
		//TODO:重写参数,取出Token和push服务
		key = k
		if key == "c" {
			c = strings.Join(v, "")
			count++
			oc = 1
		}
		if key == "w" {
			ws = strings.Join(v, "")
			wi, err = strconv.Atoi(ws)
			if err != nil {
				problem++
			}
			count++
		}
		if key == "h" {
			hs = strings.Join(v, "")
			h, err = strconv.Atoi(hs)
			if err != nil {
				problem++
			}
			count++
		}
		_ = v
	}
	if count == 3 && problem == 0 {
		content := qrcode.ToPng(c, wi, h)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "Image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(content.Bytes())))
		_, _ = w.Write(content.Bytes())
		count = 0
		problem = 0
	} else if count == 1 && oc == 1 {
		content := qrcode.NoToPng(ws)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "Image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(content.Bytes())))
		_, _ = w.Write(content.Bytes())
		count = 0
		problem = 0
	} else {
		w.WriteHeader(403)
		_, _ = fmt.Fprintf(w, "Error pram set!")
		count = 0
		problem = 0
	}
}
func httpServerPid() {
	http.HandleFunc("/", qr)                        //处理 push
	err := http.ListenAndServe("0.0.0.0:8081", nil) //端口为8787
	if err != nil {
		log.Debug("ListenAndServe: ", err) //监听端口
	}
}

// HTTPServer http 服务器示例
type HTTPServer struct {
	http *http.Server
	cmd  *cobra.Command
}

// PidSavePath pid保存路径
func (httpServer *HTTPServer) PidSavePath() string {
	return "./"
}

// Name pid文件名
func (httpServer *HTTPServer) Name() string {
	return "http"
}

// SetCommand 从 daemon 获得 cobra.Command 对象

func (httpServer *HTTPServer) Start() {
	httpServerPid()
}
func (httpServer *HTTPServer) Stop() error {
	err := httpServer.http.Shutdown(context.Background())
	return err
}

// Restart 重启web服务前关闭http服务
func (httpServer *HTTPServer) Restart() error {
	err := httpServer.Stop()
	return err
}

func main() {
	// 自定义输出文件
	out, _ := os.OpenFile("./http.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	err, _ := os.OpenFile("./http_err.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	// 初始化一个新的运行程序
	proc := daemon.NewProcess(new(HTTPServer)).SetPipeline(nil, out, err)
	proc.On(syscall.SIGTERM, func() { fmt.Println("a custom signal") })
	// 示例,多级命令服务
	// 这里的示例由于实现了 Command 接口, 所以这里会出现 flag test 不存在的情况, 实际情况, 每一个 worker 都应该是唯一的
	// 不要共享一个 worker 对象指针
	daemon.GetCommand().AddWorker(proc).AddWorker(proc)
	// 示例,主服务
	daemon.Register(proc)

	// 运行
	if rs := daemon.Run(); rs != nil {
		log.Info(rs)
	}
}
