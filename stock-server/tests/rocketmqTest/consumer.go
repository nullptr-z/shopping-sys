package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.1.107:9876"}),
		consumer.WithGroupName("stock_group"),
	)

	if err := c.Subscribe("topic1", consumer.MessageSelector{}, subCallBack); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()
	//不能让主goroutine退出
	time.Sleep(time.Hour)
	_ = c.Shutdown()
}

func subCallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	// 打印接收到的消息
	for _, msg := range msgs {
		fmt.Println("打印消息:", msg)
	}

	// 返回成功到主服务器，表示接受到了消息
	return consumer.ConsumeSuccess, nil
}
