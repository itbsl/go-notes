package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

//计算MD5的方式
func getMD5(str []byte) string {
	//1.计算数据的md5
	result := md5.Sum(str)
	//将数据src编码为16进制字符串s，每个字节占两位，不足补0。
	return hex.EncodeToString(result[:])
}

func main() {
	fmt.Println(getMD5([]byte("111111")))
}