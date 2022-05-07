package main

import (
	"flag"
	"fmt"
	"net"
	"runtime"
	"time"
)

func main() {

	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores - 1)

	domainP := flag.String("domain", "", "请输入指定的IP地址")

	flag.Parse()

	domain := *domainP

	if domain == "" {
		fmt.Println("请输入IP地址")
		return
	}

	openPort := make([]int, 10)

	exitChan := make(chan bool)

	portChan := make(chan int)

	resChan := make(chan int, 100)

	go func() {
		for i := 1; i < 65536; i++ {
			portChan <- i
		}
		close(portChan) //及时关闭管道 使用for循环取出时才不会等待
	}()
	timeout := time.Millisecond * 200

	for i := 0; i < 18; i++ {
		go CheckIsOpen(domain, timeout, portChan, exitChan, openPort, resChan)
	}

	for i := 0; i < 18; i++ {
		<-exitChan
	}

	close(exitChan)
	close(resChan)
	fmt.Println("开启的端口号:")
	for {
		port, res := <-resChan

		if !res {
			break
		}

		fmt.Println(port)

	}
}

//检查端口是否打卡
func CheckIsOpen(domain string, timeout time.Duration, portChan chan int, exitChan chan bool, openPort []int, reschan chan int) {
	for {
		port, ok := <-portChan

		if !ok {
			break
		}

		socket := fmt.Sprintf("%v:%v", domain, port)

		_, err := net.DialTimeout("tcp", socket, timeout) //发起握手操作  有回应 说明打开了端口 并且做了延时 超过一定时间 直接说明未打开

		defer func() {
			recover()
		}()

		if err == nil {
			reschan <- port
			fmt.Println(socket, " 打开")
			openPort = append(openPort, port)
		} else {
			fmt.Println(socket, " 关闭")
		}
	}

	exitChan <- true

}
