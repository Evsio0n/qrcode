package main

import (
	"./module/library/log"
	"./module/library/qrcode"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
			oc=1
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
		w.Header().Set("Content-Length",strconv.Itoa(len(content.Bytes())))
		_,_=w.Write(content.Bytes())
		count = 0
		problem = 0
	}else if count==1&&oc==1{
		content := qrcode.NoToPng(ws)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "Image/png")
		w.Header().Set("Content-Length",strconv.Itoa(len(content.Bytes())))
		_,_=w.Write(content.Bytes())
		count = 0
		problem = 0
	}else {
		w.WriteHeader(403)
		_, _ = fmt.Fprintf(w, "Error pram set!")
		count = 0
		problem = 0
	}
}
func main() {
	http.HandleFunc("/", qr)                        //处理 push
	err := http.ListenAndServe("0.0.0.0:8787", nil) //端口为8787
	if err != nil {
		log.Debug("ListenAndServe: ", err) //监听端口
	}
}
