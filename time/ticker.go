package main

import (
	"fmt"
	"time"
)

//Golang 定时器包括：一次性定时器（Timer）和周期性定时器(Ticker)。
//timer创建有两种方式，time.NewTimer(Duration) 和time.After(Duration)。 后者只是对前者的一个包装。
//它会按照一个时间间隔往channel发送系统当前时间，而channel的接收者可以以固定的时间间隔从channel中读取事件。
//Ticker 跟 Timer 的不同之处，就在于 Ticker 时间达到后不需要人为调用 Reset 方法，会自动续期

//NewTicker() 返回一个新的Ticker，该Ticker包含一个通道字段，并会每隔时间段d就向该通道发送当时的时间。
//它会调整时间间隔或者丢弃tick信息以适应反应慢的接收者。如果d<=0会panic。关闭该Ticker可以释放相关资源。
//Stop() 关闭一个Ticker。在关闭后，将不会发送更多的tick信息。
//Stop不会关闭通道t.C，以避免从该通道的读取不正确的成功
func basicTickerTutorial() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		time.Sleep(time.Second * 10)
		done <- true
	}()

	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			fmt.Printf("当前时间为: %s\n", t.Format("2006-01-02 15:04:05"))
		}
	}
}

//time.NewTicker 周期时间到了，但是之前程序没有执行完,怎么处理？
//time.NewTicker定时触发执行任务，当下一次执行到来而当前任务还没有执行结束时，会等待当前任务执行完毕后再执行下一次任务。
func advancedTutorial() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		time.Sleep(time.Second * 10)
		done <- true
	}()

	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case <-ticker.C:
			fmt.Printf("当前时间为: %s\n", time.Now().Format("2006-01-02 15:04:05"))
			time.Sleep(time.Second * 4)
		}
	}
}

func main() {
	advancedTutorial()
}
