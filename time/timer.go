package main

import (
	"fmt"
	"time"
)

//Golang 定时器包括：一次性定时器（Timer）和周期性定时器(Ticker)。
//timer创建有两种方式，time.NewTimer(Duration) 和time.After(Duration)。 后者只是对前者的一个包装。

//基础教程
func basicTimerTutorial() {
	//打印当前时间
	fmt.Printf("当前时间为: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	//创建定时器，2秒后，定时器就会向自己的C字段发送一个time.Time类型的值
	timer := time.NewTimer(time.Second * 2)
	//阻塞两秒
	fmt.Printf("当前时间为: %s\n", (<-timer.C).Format("2006-01-02 15:04:05"))
}

// reset使用教程
// func (t *Timer) Reset() bool
// 重置的动作实质上是先停止定时器，再启动，其返回值即停止定时器(Stop())的返回值。
// 需要注意的是，重置定时器虽然可以用于修改还未超时的定时器，但正确的使用方式还是
// 针对已过期的或已被停止的定时器，同时其返回值也不可靠，返回值存在的价值仅仅是与
// 前面的版本兼容。
// 实际上，重置定时器意味着通知系统守护协程移除该定时器，重新设定时间后，再把定时器
// 交给守护协程。
func resetTimerTutorial() {
	//打印当前时间
	fmt.Printf("当前时间为: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	timer := time.NewTimer(time.Second * 10)
	//重置定时器
	timer.Reset(time.Second * 2)
	//获取当前时间
	fmt.Printf("当前时间为: %s\n", (<-timer.C).Format("2006-01-02 15:04:05"))

	//timer到固定时间后会执行一次，请注意是一次，而不是多次。
	//但是可以通过reset来实现每隔固定时间段执行。使用timer定时器，超时后需要重置，才能继续触发。
	timer = time.NewTimer(time.Second * 2)
	fmt.Printf("当前时间为: %s\n", (<-timer.C).Format("2006-01-02 15:04:05"))
	//执行Reset方法后，有可以重新从Timer.C中取值了
	//基于此，我们甚至可以使用reset方法来实现每隔固定时间执行任务
	timer.Reset(time.Second * 3)
	fmt.Printf("当前时间为: %s\n", (<-timer.C).Format("2006-01-02 15:04:05"))

	timer = time.NewTimer(time.Second * 2)
	var i int
	for {
		select {
		case <-timer.C:
			i++
			fmt.Printf("count: %v\n", i)
			timer.Reset(time.Second * 2) //每次使用完后需要人为重置下
		}
	}
	timer.Stop()
}

// stop使用教程
// func (t *Timer) Stop() bool
// 其返回值代表定时器有没有超时
// true: 定时器超时前停止，后续不会再发送事件；
// false: 定时器超时后停止。
// 实际上，停止定时器意味着通知系统守护协程移除该定时器
func stopTimerTutorial() {
	//打印当前时间
	fmt.Printf("当前时间为: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	timer := time.NewTimer(time.Second * 2)
	isStopped := timer.Stop()
	fmt.Printf("停止结果为: %t\n", isStopped)

	timer = time.NewTimer(time.Second * 2)
	fmt.Printf("当前时间为: %s\n", (<-timer.C).Format("2006-01-02 15:04:05"))
	isStopped = timer.Stop()
	fmt.Printf("停止结果为: %t\n", isStopped)

	timer = time.NewTimer(time.Second * 2)
	go func() {
		<-timer.C
		fmt.Println("异步任务等待2秒执行")
	}()
	isStopped = timer.Stop()
	if isStopped {
		fmt.Println("定时器被停止了")
	}

	timer = time.NewTimer(time.Second * 2)
	go func() {
		<-timer.C
		fmt.Println("异步任务等待2秒执行")
	}()
	time.Sleep(time.Second * 3) //sleep 3秒使得异步任务执行成功，也就是 <-timer.C
	isStopped = timer.Stop()
	if !isStopped {
		fmt.Println("定时器停止失败")
	}
}

//after使用教程
func afterTimerTutorial() {
	//打印当前时间
	fmt.Printf("当前时间为: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	//有时我们就是想等待指定的时间，没有提前停止定时器的需求，也没有复用该定时器的需求，那么可以使用匿名的定时器
	fmt.Printf("当前时间为: %s\n", (<-time.After(time.Second * 2)).Format("2006-01-02 15:04:05"))

	//AfterFunc()函数用于等待经过的时间，此后，它将在其自己的go-routine中调用已定义的函数“f”。
	//此外，此函数在时间包下定义。在这里，您需要导入“time”包才能使用这些函数。
	var c = make(chan struct{}, 1)
	time.AfterFunc(time.Second*2, func() {
		fmt.Printf("执行业务，当前时间为: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		c <- struct{}{}
	})
	<-c
}

func main() {
	//basicTimerTutorial()
	//resetTimerTutorial()
	//stopTimerTutorial()
	afterTimerTutorial()
}
