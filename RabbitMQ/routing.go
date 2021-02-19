package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

//创建路由模式下的RabbitMQ实例
func NewRoutingRabbitMQ(url string, exchange string, routingKey string) *rabbitMQ {
	return newRabbitMQ(url, "", exchange, routingKey, "route")
}

func (this *rabbitMQ) routeSend(message string) {
	//声明交换机
	err := this.channel.ExchangeDeclare(
		this.exchange,
		"direct",
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
//路由模式：消费者客户端
func (this *rabbitMQ) routeConsume(handleFunc HandleFunc) {
	//声明交换机
	err := this.channel.ExchangeDeclare(
		this.exchange,
		"direct",
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
