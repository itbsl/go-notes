package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"strconv"
)

package main

import (
"fmt"
"gopkg.in/gomail.v2"
"strconv"
)

//发送邮件
func SendMail(mailTo []string, subject string, body string) error {

	//定义邮箱服务器连接信息，如果是阿里云邮箱，password填密码，腾讯邮箱填授权码
	mailConfig := map[string]string{
		"host":     "smtp.exmail.qq.com",
		"username": "bashlog@qq.cn",
		"password": "授权码",
		"port":     "465",
	}
	port, _ := strconv.Atoi(mailConfig["port"])

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", mailConfig["username"], "白兮")
	msg.SetHeader("To", mailTo...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	msg.Attach("/Users/itbsl/Desktop/WechatIMG1.jpeg")

	dialer := gomail.NewDialer(mailConfig["host"], port, mailConfig["username"], mailConfig["password"])
	err := dialer.DialAndSend(msg)

	return err
}

func main() {

	//收件人
	mailTo := []string{
		"1xxxx@foxmail.com",
		"2xxxx@foxmail.com",
	}

	//邮件主题
	subject := "这是主题"

	//邮件正文
	body := "这是邮件正文"

	err := SendMail(mailTo, subject, body)
	if err != nil {
		fmt.Println(err)
	}
}
