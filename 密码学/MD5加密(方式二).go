package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

//计算数据的MD5值
func getMD5(str string) string {
	//1.创建一个使用MD5校验的hash.Hash接口
	hash := md5.New()

	//2.通过io操作将数据写入hash对象中
	io.WriteString(hash, str)

	//3.计算结果
	result := hash.Sum(nil)

	//4.将数据src编码为16进制字符串s，每个字节占两位，不足补0，最终返回一个32位长度的字符串
	//return fmt.Sprintf("%x", result)
	return hex.EncodeToString(result)
}

func main() {
	fmt.Println(getMD5("111111"))
}