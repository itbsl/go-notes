package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

//爬取数据
//爬取《王者荣耀》
//第一页：https://tieba.baidu.com/f?kw=%E7%8E%8B%E8%80%85%E8%8D%A3%E8%80%80&ie=utf-8&cid=&tab=corearea&pn=0
//第二页：https://tieba.baidu.com/f?kw=%E7%8E%8B%E8%80%85%E8%8D%A3%E8%80%80&ie=utf-8&cid=&tab=corearea&pn=50
//第三页：https://tieba.baidu.com/f?kw=%E7%8E%8B%E8%80%85%E8%8D%A3%E8%80%80&ie=utf-8&cid=&tab=corearea&pn=100
func spider(index int) {
	defer wg.Done()
	url := "https://tieba.baidu.com/f?kw=%E7%8E%8B%E8%80%85%E8%8D%A3%E8%80%80&ie=utf-8&cid=&tab=corearea&pn="+strconv.Itoa(50 * (index - 1))
	result, err := httpGet(url)
	if err != nil {
		fmt.Printf("http.Get()函数执行错误,错误为:%v\n", err)
		return
	}
	//将抓取到的内容写到文件中
	file, err := os.OpenFile("index"+strconv.Itoa(index)+".html", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Printf("os.OpenFile()函数执行错误,错误为:%v\n", err)
		return
	}
	//关闭文件
	defer file.Close()

	file.WriteString(result)
}

func httpGet(url string) (result string, err error) {
	resp, err := http.Get(url)
	if err != nil {//若http请求发生错误，直接返回给调用者
		return
	}
	defer resp.Body.Close()

	//循环读取网页数据，返回给调用者
	buf := make([]byte, 4096)
	for {
		//注意：for循环里的err是一个新的变量，不同于函数返回值里定义的err(因为这是一个新的块作用域)
		//如果怕混淆，那么就定义个其它的变量名称，比如err1,err2
		n, err := resp.Body.Read(buf)
		//累加每一次循环读取到哦buf数据，存入result
		result += string(buf[:n])
		if err != nil {
			if err == io.EOF { //读取到末尾
				fmt.Printf("读取网页完成\n")
				break
			} else {
				return result, err
			}
		}
	}
	return
}

func main() {
	//步骤：
	//1.明确目标URL
	//2.发送请求，获取应答数据包
	//3.保存、过滤数据。提取有用信息
	//4.使用、分析得到数据信息
	//指定爬取起始、终止页
	var start, end int
	fmt.Print("请输入爬取的起始页(>=1):")
	fmt.Scanf("%d", &start)

	fmt.Printf("请输入爬取的终止页(>=%d):", start)
	fmt.Scanf("%d", &end)

	for i := start; i <= end; i++ {
		wg.Add(1)
		go spider(i)
	}
	wg.Wait()
}