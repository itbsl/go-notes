package main

import (
	"fmt"
	"time"
)

//select是一种go可以处理多个通道之间的机制，看起来和switch语句很相似，
//但是select其实和IO机制中的select一样，多路复用通道，随机选取一个进行执行，
//如果说通道(channel)实现了多个goroutine之前的同步或者通信，
//那么select则实现了多个通道(channel)的同步或者通信，并且select具有阻塞的特性。

//select 是 Go 中的一个控制结构，类似于用于通信的 switch 语句。每个 case 必须是一个通信操作，要么是发送要么是接收。

//select 随机执行一个可运行的 case。如果没有 case 可运行，它将阻塞，直到有 case 可运行。一个默认的子句应该总是可运行的。

//golang 的 select 就是监听 IO 操作，当 IO 操作发生时，
//触发相应的动作每个case语句里必须是一个IO操作，确切的说，应该是一个面向channel的IO操作。

//如果select中的信道都阻塞的话，就会立即进入 default 分支，并不会阻塞。
//但是如果没有 default 语句，则会阻塞直到某个信道操作成功为止。

/*
 * select语句只能用于信道的读写操作
 * 对于case条件语句中，如果存在信道值为nil的读写操作，则该分支将被忽略，可以理解为从select语句中删除了这个case语句
 * 如果有超时条件语句，判断逻辑为如果在这个时间段内一直没有满足条件的case，则执行这个超时case。
   如果此段时间内出现了可操作的case,则直接执行这个case。一般用超时语句代替了default语句。
 * 对于空的select{}，会引起死锁
 * 对于for中的select{}, 也有可能会引起cpu占用过高的问题
*/

func basicTutorial() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		time.Sleep(time.Minute * 60)
		ch1 <- 1
	}()
	go func() {
		time.Sleep(time.Second * 60)
		ch2 <- 2
	}()

	for {
		select {
		case i := <-ch1:
			fmt.Printf("从ch1读取的内容为: %v\n", i)
		case j := <-ch2:
			fmt.Printf("从ch2读取的内容为: %v\n", j)
		default:
			fmt.Printf("当前时间为: %v\n", time.Now().Format("2006-01-02 15:04:05"))
			time.Sleep(time.Second)
		}
	}
}

//竞争选举
func raceTutorial() {
	c1, c2, c3 := make(chan int), make(chan int), make(chan int)
	go func() {
		c1 <- 1
	}()
	go func() {
		c2 <- 2
	}()
	go func() {
		c3 <- 3
	}()
	select {
	case i := <-c1:
		fmt.Printf("从c1读取了数据: %v\n", i)
	case j := <-c2:
		fmt.Printf("从c2读取了数据: %v\n", j)
	case k := <-c3:
		fmt.Printf("从c3读取了数据: %v\n", k)
	}
}

// 超时处理，保证不阻塞
func timeoutTutorial() {
	var c = make(chan int)
	go func() {
		time.Sleep(time.Second * 10)
		c <- 1
	}()
	select {
	case v := <-c:
		fmt.Printf("获取的数据为: %v\n", v)
	case <-time.After(time.Second * 5):
		fmt.Println("timeout!!!")
	}
}

func isBufferBlockedTutorial() {
	var ch = make(chan int, 5)
	go func() {
		//time.Sleep(time.Second)
		for {
			<-ch
			time.Sleep(time.Second * 5)
		}
	}()

	for {
		select {
		case ch <- 1:
			fmt.Println("添加成功！！！")
			time.Sleep(time.Second)
		default:
			fmt.Println("资源已满，请稍后再试")
			time.Sleep(time.Second)
		}
	}
}

func main() {
	//basicTutorial()
	//raceTutorial()
	//timeoutTutorial()
	isBufferBlockedTutorial()
	//阻塞main函数
	//有时候我们会让main函数阻塞不退出，如http服务，我们会使用空的select{}来阻塞main goroutine
	//这行主函数就永远阻塞住了，这里要注意一定要有一只活动的goroutine，否则会报deadlock。
	//可以把select{}换成for{}试一下，打开系统管理器看下CPU占用变化
	//var ch = make(chan int)
	//go func() {
	//	for {
	//		ch <- 1
	//		time.Sleep(time.Second)
	//	}
	//}()
	//go func() {
	//	for {
	//		fmt.Println(<-ch)
	//	}
	//}()
	//select {}
}
