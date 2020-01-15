package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

//计算数据的MD5值
func getMD5(data []byte) string {
	//1.计算数据data的MD5校验和
	result := md5.Sum(data)
	//2.将数据src编码为16进制字符串s，每个字节占两位，不足补0，最终返回一个32位长度的字符串
	return hex.EncodeToString(result[:])
}

func main() {
	fmt.Println(getMD5([]byte("111111")))
}