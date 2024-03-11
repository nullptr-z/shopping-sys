package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	producer, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.1.107:9876"}))
	if err != nil {
		panic(err)
	}

	err = producer.Start()
	if err != nil {
		panic(fmt.Sprint("start producer error", err))
	}

	resp, err := producer.SendSync(context.Background(), &primitive.Message{
		Topic: "topic1",
		Body:  []byte("hi topic1"),
	})
	if err != nil {
		panic(fmt.Sprint("send message error", err))
	} else {
		fmt.Println("send message success: ", resp.String())
	}

	if err = producer.Shutdown(); err != nil {
		panic(fmt.Sprint("shutdown producer error", err))
	}
	fmt.Println("rocketmq shutdown")
}
