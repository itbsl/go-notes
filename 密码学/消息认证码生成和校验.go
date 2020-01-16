package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

//生成消息认证码
func GenerateHmac(plaintext, key []byte) []byte {
	//1.创建一个采用hash.Hash作为底层hash接口、key作为密钥的HMAC算法的hash接口
	//需要指定使用的哈希算法和秘钥
	hash := hmac.New(sha256.New, key)

	//2.给hash对象添加数据
	hash.Write(plaintext)

	//3.计算散列值
	hashText := hash.Sum(nil)

	return hashText
}

//验证消息认证码
func VerifyHmac(plaintext, key, hashText []byte) bool {
	//1.创建一个采用hash.Hash作为底层hash接口、key作为密钥的HMAC算法的hash接口
	//需要指定使用的哈希算法和秘钥
	hash := hmac.New(sha256.New, key)

	//2.给hash对象添加数据
	hash.Write(plaintext)

	//3.计算散列值
	mac := hash.Sum(nil)

	//4.比较生成的散列值和参数中的散列值hashText
	return hmac.Equal(hashText, mac)
}

func main() {

	key := []byte("12345678")
	plaintext := []byte("我们结婚吧!")
	mac := GenerateHmac(plaintext, key)
	fmt.Println(VerifyHmac(plaintext, key, mac))
}
