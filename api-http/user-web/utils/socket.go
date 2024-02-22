package utils

import (
	"net"
)

// 随机得到一个空闲端口
func GetFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0
	}
	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0
	}
	return lis.Addr().(*net.TCPAddr).Port
}
