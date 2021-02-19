package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

func NewTopicRabbitMQ(url string, exchange string, routingKey string) *rabbitMQ {
	return newRabbitMQ(url, "", exchange, routingKey, "topic")
}

func (this *rabbitMQ) topicSend(message string) {
	//声明交换机
	err := this.channel.ExchangeDeclare(
		this.exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
		)
	this.failOnError(err, "交换机创建失败")

	//发送消息
	err = this.channel.Publish(
		this.exchange,
		this.key, //要设置
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	this.failOnError(err, "消息发送失败")
}

func (this *rabbitMQ) topicConsume(handleFunc HandleFunc) {
	//声明交换机
	err := this.channel.ExchangeDeclare(
		this.exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	this.failOnError(err, "交换机创建失败")

	//创建队列
	queue, err := this.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
		)
	this.failOnError(err, "创建队列失败")

	//绑定队列到exchange中
	err = this.channel.QueueBind(
		queue.Name,
		this.key,
		this.exchange,
		false,
		nil,
		)
	this.failOnError(err, "队列绑定失败！")

	//注册消费者
	messages, err := this.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	this.failOnError(err, "注册消费者失败！")

	forever := make(chan bool)
	go func() {
		for data := range messages {
			//实现我们要处理的逻辑
			handleFunc(string(data.Body))
		}
	}()
	fmt.Printf("[*] 等待消息，使用Ctrl+C退出\n")
	<-forever
}