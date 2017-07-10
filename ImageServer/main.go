// ImageServer project main.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	Template_Dir = "./views/"
	Upload_Dir   = "./upload/"
)

//解析config.json 接收使用
type Config struct {
	ListenAddr string
	Storage    string
}

//返回消息使用的
type mesage struct {
	Code int
	Msg  string
}

//全局声明
var conf Config

//载入config.json
func LoadConf() {
	r, e := os.Open("etc/config.json")
	if e != nil {
		log.Fatal(e)
	}
	defer r.Close()
	//解析json
	decode := json.NewDecoder(r)
	//解析到Config自定义类型
	err := decode.Decode(&conf)
	if err != nil {
		log.Fatal(err)
	}

}

//上传方法
func UploadHandler(w http.ResponseWriter, req *http.Request) {
	//判断请求方法 post
	if req.Method != "POST" {
		/**  方法一 直接处理
		//初始化返回消息
		msg := mesage{Code: 0, Msg: "Don't support GET!!!"}

		//封装返回消息json
		msg1, e := json.Marshal(msg)
		if e != nil {

		}
		//返回错误码
		w.WriteHeader(http.StatusMethodNotAllowed)
		//返回消息
		w.Write(msg1)
		**/

		//方法二 调用模板 定向到上传页面
		t, _ := template.ParseFiles(Template_Dir + "file.html")
		t.Execute(w, "上传文件")

	} else {
		req.ParseMultipartForm(32 << 20)

		file, handler, err := req.FormFile("uploadfile")
		//获取上传文件扩展名
		fileext := filepath.Ext(handler.Filename)
		//生成新文件名
		filename := strconv.FormatInt(time.Now().Unix(), 10) + fileext
		f, _ := os.OpenFile(Upload_Dir+filename, os.O_CREATE|os.O_WRONLY, 0660)
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			fmt.Fprintf(w, "%v", "上传失败")
			return
		}
		filedir, _ := filepath.Abs(Upload_Dir + filename)
		fmt.Fprintf(w, "%v", filename+"上传完成,服务器地址:"+filedir)

	}

	//获取请求参数必须执行方法
	req.ParseForm()
	if len(req.Form) > 0 {
		for k, v := range req.Form {
			fmt.Println(k, v)

		}

	}

}

//下载方法
func DownLoadHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if len(req.Form) > 0 {
		imgid := req.Form.Get("id")
		if len(imgid) > 0 {
			fmt.Println("imgID:", imgid)
			fmt.Println(req.Method)
			os.IsExist()
			http.ServeFile(w, req, Upload_Dir+imgid)
		}

	}

}

func main() {
	//载入config.json
	LoadConf()
	//路由请求
	http.HandleFunc("/up", UploadHandler)
	http.HandleFunc("/dw/", DownLoadHandler)
	//监听
	http.ListenAndServe(conf.ListenAddr, nil)

}
