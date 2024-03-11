package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type OrderListener struct{}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地逻辑")
	time.Sleep(time.Second * 2)
	fmt.Println("执行本地逻辑失败")
	//本地执行逻辑无缘无故失败 代码异常 宕机
	return primitive.UnknowState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("rocketmq的消息回查")
	time.Sleep(time.Second * 2)
	return primitive.CommitMessageState
}

func main() {
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{},
		producer.WithNameServer([]string{"192.168.1.107:9876"}),
	)
	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	res, err := p.SendMessageInTransaction(
		context.Background(),
		primitive.NewMessage("TransTopic", []byte("ttttttttt1t222")),
	)
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", res.String())
	}

	time.Sleep(10 * time.Second)
	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败")
	}
}
