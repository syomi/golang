// myHttpServer project main.go
package main

import (
	"log"
	"net/http"
	"os"
)

func myhandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ip := r.FormValue("ip")
	if len(ip) < 1 {
		w.Write([]byte("请输入?ip=IP"))
	}
	er := writeFile(ip)
	if er != nil {
		w.Write([]byte("保存失败"))
	} else {
		w.Write([]byte("保存成功"))
	}

}

func writeFile(content string) (er error) {
	file, err := os.OpenFile("save.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, e := file.WriteString(content + "\n")

	return e
}

func main() {

	http.HandleFunc("/report", myhandler)
	http.ListenAndServe(":9999", nil)
}
