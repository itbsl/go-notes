package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

//计算数据的Sha256值
func getSha256(src string) string {
	//1.打开文件
	file, err := os.Open(src)
	if err != nil {
		log.Fatalf("os.Open(src)函数执行错误,错误为:%v\n", err)
	}
	defer file.Close()

	//2.创建基于sha256算法的Hash对象
	hash := sha256.New()

	//3.将文件数据拷贝给hash对象
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Fatalf("io.Copy()函数执行错误,错误为:%v\n", err)
	}

	//4.计算文件的哈希值
	result := hash.Sum(nil)

	//5.将数据src编码为16进制字符串s，每个字节占两位，不足补0，最终返回一个64位长度的字符串
	return hex.EncodeToString(result)
}

func main() {
	fmt.Println(getSha256("private.pem"))
}
