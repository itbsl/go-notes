package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

//生成RSA秘钥，并且保存到磁盘文件中
func GenerateRSAKey(keySize int) {
	//1.使用RSA中的GenerateKey方法生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		log.Fatalf("rsa.GenerateKey()函数执行错误,错误为:%v\n", err)
	}

	//2.通过x509标准将得到的RSA私钥序列化为ASN.1 PKCS#1 DER编码
	derText := x509.MarshalPKCS1PrivateKey(privateKey)

	//3.初始化一个pem.Block(为了实现PEM数据编码)
	block := pem.Block{
		Type:    "RSA PRIVATE KEY", //得自前言的类型（如"RSA PRIVATE KEY"）
		Headers: nil,               //可选的头项
		Bytes:   derText,           //内容解码后的数据，一般是DER编码的ASN.1结构
	}
	//注:最终生成的文件会以"-----BEGIN RSA PRIVATE KEY-----"开头，以"-----END RSA PRIVATE KEY-----"结束

	//4.pem编码
	//4.1创建一个文件用于保存生成的私钥
	file, err := os.Create("private.pem")
	if err != nil {
		log.Fatalf("os.Create()函数执行错误，错误为:%v\n", err)
	}
	defer file.Close()
	//4.2对数据进行PEM编码
	err = pem.Encode(file, &block)
	if err != nil {
		log.Fatalf("pem.Encode()函数执行错误，错误为:%v\n", err)
	}
}

func main() {
	//生成RSA秘钥
	GenerateRSAKey(1024)
}
