package main

import (
	"fmt"
	"net/rpc"
	"time"
	"net"
	"strings"
)

func getLocalIP() (ip string) {
	conn, err := net.Dial("udp", "www.baidu.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	ip = strings.Split(conn.LocalAddr().String(), ":")[0]
	fmt.Println("本地IP：", ip)
	return ip
}

func main() {
	clint, err := rpc.DialHTTP("tcp", "114.212.87.65:1234")
	if err != nil {
		fmt.Println("服务出错0：", err)
	}
	var serverTime, serverTime0, localTime, loaclTime0 time.Time
	var minSigma, sigma, localTimeDiff, serverTimeDiff time.Duration
	args := getLocalIP()
	minSigma = 100000

	for i := 1; i<=100; i++ {
		//fmt.Println(i)
		localTime = time.Now()
		err = clint.Call("Watcher.GetServerTime", args, &serverTime)
		if err != nil {
			fmt.Println("服务出错1：", err)
		}
		if i == 1 {
			//fmt.Println(localTime)
			serverTime0 = serverTime
			loaclTime0 = localTime
		}

		sigma = time.Since(localTime)
		if minSigma > sigma {
			minSigma = sigma
		}
	}

	localTimeDiff = localTime.Sub(loaclTime0)
	serverTimeDiff = serverTime.Sub(serverTime0)

	//fmt.Println("localTimeDiff : ", localTimeDiff)
	//fmt.Println("serverTimeDiff : ", serverTimeDiff)
	//fmt.Println("sigma : ", localTimeDiff-serverTimeDiff)

	minSigma -= serverTimeDiff-localTimeDiff

	fmt.Println("serverTime : ", serverTime0)
	fmt.Println("sigma : ", minSigma/2)
	fmt.Println("real time from server:", serverTime0.Add(minSigma/2))

}