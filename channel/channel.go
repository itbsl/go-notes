package main

import (
	"fmt"
	"sync"
	"time"
)

//管道实现了一种FIFO(先入先出)的队列，数据总是按照写入的顺序流出管道。
//内置函数`len()`和`cap()`作用于管道，分别用于查询缓冲区中数据的个数及缓冲区的大小。

// Go语言中 nil 的地址为 0x0
// Go语言中 chan 为引用类型，可以与nil比较是否相等
// 读nil管道会阻塞，但是如果只有读nil管道的操作，而没有写nil管道的操作会发生死锁
// 写nil管道会阻塞，但是如果只有写nil管道的操作，而没有读nil管道的操作会发生死锁
func nilChanTutorial() {
	var c chan int
	if c == nil {
		fmt.Printf("c的值为: %v, %p\n", c, c)
	}
	//关闭 nil 管道会引发panic
	//close(c)
}

var wg sync.WaitGroup

//`管道没有缓冲区时(例如：c = make(chan int))，从管道读取数据时会阻塞，直到有协程向管道中写入数据。
//类似地，向管道写入数据数据也会阻塞，直到有协程从管道读取数据。
func noBufferedChanTutorial() {
	var c = make(chan int)
	wg.Add(2)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 5)
		c <- 1
		fmt.Println("写入了值")
	}()
	go func() {
		defer wg.Done()
		<-c
		fmt.Println("获取了值")
	}()
	wg.Wait()
	//无缓冲管道可以被关闭

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			c <- i
			time.Sleep(time.Second)
		}
		//写入完成后，记得关闭管道,否则range 时会panic
		close(c)
	}()
	go func() {
		defer wg.Done()
		for v := range c {
			fmt.Printf("读取到的值为: %v\n", v)
		}
	}()
	wg.Wait()
}

func bufferedChanTutorial() {
	var c = make(chan int, 5)
	c <- 1
	c <- 2
	c <- 3
	fmt.Printf("读取到的内容为: %v\n", <-c)
}

//关闭通道使用教程
func closeChanTutorial() {
	//nil管道不可被关闭，使用close关闭nil管道会引发panic
	var c1 chan int
	//close(c1) //panic
	fmt.Printf("nil管道的长度为: %v, 容量为: %v\n", len(c1), cap(c1))

	//无缓冲区的管道可以被close
	var c2 = make(chan int)
	fmt.Printf("无缓冲区的管道长度为: %v, 容量为: %v\n", len(c2), cap(c2))
	close(c2)
	//无缓冲管道关闭后，读取的值都是对应类型的零值
	fmt.Printf("%v, %v\n", <-c2, <-c2) //输出0，0
	for v := range c2 {
		fmt.Printf("读取到的值为: %v\n", v)
	}

	//有缓冲区的管道可以被close
	var c3 = make(chan int, 2)
	fmt.Printf("有缓冲区的管道长度为: %v, 容量为: %v\n", len(c3), cap(c3))
	c3 <- 1
	c3 <- 2
	//如果一个有缓冲管道没有关闭，且没有写操作，使用range读会导致panic
	//如果一个有缓冲管道关闭了，即使没有写操作，使用range也可以读取而不导致panic
	//使用内置函数close()可以关闭管道，尝试向关闭的管道写入数据会触发panic，但关闭的管道仍可读。
	close(c3)
	for v := range c3 {
		fmt.Printf("获取到的值为: %v\n", v)
	}

	//关闭已经被关闭的管道会导致panic
	var c4 = make(chan int)
	close(c4)
	close(c4) //会导致panic
}

func readFromChanTutorial() {
	var c = make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	//管道读取表达式可以给两个变量赋值
	value := <-c
	fmt.Printf("读取到的值为: %v\n", value)
	//第一个变量表示读出的数据，第二个变量(bool类型)表示是否成功读取了数据，
	//需要注意的是，第二个变量不用于指示管道的关闭状态。
	value, ok := <-c
	fmt.Printf("读取到的值为: %v, %v\n", value, ok)
	value, ok = <-c
	fmt.Printf("读取到的值为: %v, %v\n", value, ok)
}

func main() {
	//nilChanTutorial()
	//noBufferedChanTutorial()
	//bufferedChanTutorial()
	//closeChanTutorial()
	//readFromChanTutorial()
}
