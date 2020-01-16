package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

//使用RSA进行数字签名
//步骤:
//1.使用RSA生成秘钥对
//2.(服务器)使用私钥进行数字签名
//3.(客户端)使用公钥进行签名认证

//RSA签名
func RSASign(plaintext []byte, privateKeyFile string) []byte {
	//1.打开磁盘的私钥文件
	file, err := os.Open(privateKeyFile)
	if err != nil {
		log.Fatalf("os.Open()函数执行出错,错误为:%v\n", err)
	}
	defer file.Close()

	//2.将私钥文件中的内容读出
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("file.Stat()方法执行出错,错误为:%v\n", err)
	}
	buf := make([]byte, fileInfo.Size())
	file.Read(buf)

	//3.使用pem对数据解码,得到了pem.Block结构体变量
	block, _ := pem.Decode(buf)

	//4.x509将数据解析成私钥结构体 -> 得到私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("x509.ParsePKCS1PrivateKey()函数执行出错,错误为:%v\n", err)
	}

	//5.创建一个哈希对象 例如:md5、sha256
	hash := sha256.New()

	//6.给哈希对象添加数据
	hash.Write(plaintext)

	//7.计算哈希值
	hashText := hash.Sum(nil)

	//8.使用RSA中的函数对散列值签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashText)
	if err != nil {
		log.Fatalf("rsa.SignPKCS1v15()函数执行出错,错误为:%v\n", err)
	}
	return signature
}

//RSA签名验证
func VerifyRSASignature(plaintext []byte, signature []byte, publicKeyFile string) bool {
	//1.打开磁盘的公钥文件,将文件内容读出 - []byte
	file, err := os.Open(publicKeyFile)
	if err != nil {
		log.Fatalf("os.Open()函数执行出错,错误为:%v\n", err)
	}
	defer file.Close()

	//2.将私钥文件中的内容读出
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("file.Stat()方法执行出错,错误为:%v\n", err)
	}
	buf := make([]byte, fileInfo.Size())
	file.Read(buf)

	//2.使用pem解码 -> 得到pem.Block结构体变量
	block, _ := pem.Decode(buf)

	//3.使用x509对pem.Block的Bytes变量中的数据进行解析 -> 得到接口
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("x509.ParsePKCS1PublicKey()函数执行出错,错误为:%v\n", err)
	}

	//4.进行类型断言 -> 得到公钥结构体
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		log.Fatalln("断言失败")
	}

	//5.对原始消息进行哈希运算(和签名使用的哈希算法一致) -> 散列值
	hashText := sha256.Sum256(plaintext)

	//签名认证 - rsa中的函数
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashText[:], signature)
	if err != nil {
		return false
	}
	return true
}

func main() {
	plaintext := []byte("111111")
	signature := RSASign(plaintext, "private.pem")
	fmt.Println(VerifyRSASignature(plaintext, signature, "public.pem"))
}
