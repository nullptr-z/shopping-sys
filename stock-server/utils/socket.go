package utils

import (
	"net"
)

// 随机得到一个空闲端口, 设置为0
func GetFreePort() int {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
