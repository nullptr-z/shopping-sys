package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

func main() {
	// 创建Consul客户端连接
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "192.168.3.9:8500" // Consul服务的地址
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
	}
	deRegister(client)
	register(client)
	allServer(client)
	filterServer(client)
}

func allServer(client *api.Client) {
	servers, err := client.Agent().Services()
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	for _, s := range servers {
		fmt.Println("s:", s.ID)
	}
}

func filterServer(client *api.Client) {
	servers, err := client.Agent().ServicesWithFilter(`Service=="user-web"`)
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	for _, s := range servers {
		fmt.Println("s:", s.ID)
	}

}

func register(client *api.Client) {
	// 添加健康检查的微服务信息
	check := &api.AgentServiceCheck{
		// 使用/health作为健康检查的端点是一种约定俗成的做法; 一定要使用以太网 IP
		HTTP:                           "http://192.168.3.9:11001/health", // 健康检查的地址
		Interval:                       "10s",                             // 健康检查的间隔时间
		Timeout:                        "2s",                              // 健康检查的超时时间
		DeregisterCriticalServiceAfter: "5s",
	}

	// 设置要注册的服务的信息
	registration := new(api.AgentServiceRegistration)
	registration.ID = "user-web"                                             // 服务ID，唯一
	registration.Name = "user-web"                                           // 服务名称
	registration.Port = 11001                                                // 服务端口
	registration.Tags = []string{"user", "login", "register", "web", "http"} // 可选标签
	registration.Address = "192.168.3.9"                                     // 服务地址
	registration.Check = check                                               // 如果不填写，默认健康的

	// 注册服务到Consul
	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("服务注册成功！")
}

func deRegister(client *api.Client) {
	err := client.Agent().ServiceDeregister("user-web")
	if err != nil {
		// log.Fatal(err)
		log.Println(err.Error())
		return
	}
	log.Println("服务注销成功！")
}
