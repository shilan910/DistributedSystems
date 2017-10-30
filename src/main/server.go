package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"strings"
	"net/rpc"
	"net"
	"net/http"
)

type Watcher int

var validateIPs map[string]string
var ipNow = ""
var cnt = 0

func (w *Watcher) GetServerTime(arg string, localTime *time.Time) error {
	if ipNow != arg {
		fmt.Println("\n收到ip为",arg,"的client的请求")
		ipNow = arg
	}
	if _, v := validateIPs[arg]; v {
		*localTime = time.Now()
		cnt ++
	}else{
		ipNow = ""
		cnt = 0
	}
	if cnt == 1 {
		fmt.Println("当前时间为：\n", *localTime)
	}

	if cnt == 100 {
		ipNow = ""
		cnt = 0
	}
	return nil
}

func readFile(filePath string) () {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	ips := strings.Split(string(data), "\n")
	validateIPs = make(map[string]string)
	for _, v := range ips[1:] {
		validateIPs[v] = "1"
	}
}

func main() {

	readFile("config")

	watcher := new(Watcher)
	rpc.Register(watcher)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("监听失败，端口可能已经被占用")
	}
	fmt.Println("server正在监听1234端口")
	http.Serve(l, nil)

}
