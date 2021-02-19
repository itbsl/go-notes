package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

//创建发布/订阅模式下的RabbitMQ实例
func NewPubSubRabbitMQ(url string, exchange string) *rabbitMQ {
	return newRabbitMQ(url, "", exchange, "", "pub_sub")
}

//发布订阅模式：生产者客户端
func (this *rabbitMQ) pubSubSend(message string) {

	//声明交换机
	err := this.channel.ExchangeDeclare(
		this.exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
		)
	this.failOnError(err, "声明交换机失败！")

	//发送消息
	err = this.channel.Publish(
		this.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	this.failOnError(err, "消息发送失败")
}

//发布订阅模式：消费者客户端
func (this *rabbitMQ) pubSubConsume(handleFunc HandleFunc) {
	//声明交换机
	err := this.channel.ExchangeDeclare(
		this.exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	this.failOnError(err, "声明交换机失败！")

	queue, err := this.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
		)
	this.failOnError(err, "声明队列失败")

	err = this.channel.QueueBind(
		queue.Name,
		"", //在pub/sub模式下，这里的key要为空
		this.exchange,
		false,
		nil,
		)

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