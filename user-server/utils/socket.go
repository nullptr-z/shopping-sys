package utils

import (
	"fmt"
	"net"
)

// 随机得到一个空闲端口, 设置为0
func GetFreePort() int {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0
	}
	defer listener.Close()
	fmt.Println("使用的端口号：", listener.Addr().(*net.TCPAddr).Port)
	return listener.Addr().(*net.TCPAddr).Port
}
