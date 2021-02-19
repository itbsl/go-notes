package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type rabbitMQ struct {
	//连接
	conn      *amqp.Connection
	//通道
	channel   *amqp.Channel
	//队列名称
	queueName string
	//交换机
	exchange  string
	//RoutingKey/BindingKey
	key       string
	//交换器类型
	exchangeType string
}

//消息处理函数
type HandleFunc func(message string)

//错误处理函数
func (this *rabbitMQ) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s，错误信息为：%v\n", msg, err)
	}
}

//创建RabbitMQ实例
func newRabbitMQ(url string, queueName string, exchange string, key string, exchangeType string) *rabbitMQ {
	mq := &rabbitMQ{
		queueName: queueName,
		exchange:  exchange,
		key:       key,
		exchangeType: exchangeType,
	}
	var err error
	//创建连接
	mq.conn, err = amqp.Dial(url)
	mq.failOnError(err, "连接RabbitMQ失败！")

	//创建通道
	mq.channel, err = mq.conn.Channel()
	mq.failOnError(err, "创建Channel失败！")

	return mq
}

//断开channel和connection
func (this *rabbitMQ) Close() {
	this.channel.Close()
	this.conn.Close()
}

//生产者客户端
func (this *rabbitMQ) Send(message string) {
	switch this.exchangeType {
	case "simple":
		this.simpleSend(message)
	case "pub_sub":
		this.pubSubSend(message)
	case "route":
		this.routeSend(message)
	case "topic":
		this.topicSend(message)
	}
}

//消费者客户端
func (this *rabbitMQ) Consume(handleFunc HandleFunc) {
	switch this.exchangeType {
	case "simple":
		this.simpleConsume(handleFunc)
	case "pub_sub":
		this.pubSubConsume(handleFunc)
	case "route":
		this.routeConsume(handleFunc)
	case "topic":
		this.topicConsume(handleFunc)
	}
}