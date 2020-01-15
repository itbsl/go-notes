package main

import (
	"crypto/sha256"
	"fmt"
)

//计算数据的Sha256值
func getSha256(data []byte) string {

	//1.计算数据的SHA256校验和
	//2.将数据的校验和按照16进制格式化输出，每个字节格式成两位，不足2位补0，最终返回一个64位长度的字符串
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func main() {
	fmt.Println(getSha256([]byte("this is a sha256 test.")))
}