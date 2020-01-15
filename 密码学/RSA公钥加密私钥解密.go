package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

//RSA公钥加密
//注意：即使是加密相同的数据，但是每次加密产生的密文也不是相同的
func RSAPublicEncrypt(plaintext []byte, publicKeyFile string) []byte {
	//1.打开文件并且读出文件内容
	file, err := os.Open(publicKeyFile)
	if err != nil {
		log.Fatalf("os.Open()函数执行出错,错误为:%v\n", err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("file.Stat()方法执行出错,错误为:%v\n", err)
	}
	buf := make([]byte, fileInfo.Size())
	_, err = file.Read(buf)

	//2.pem解码
	block, _ := pem.Decode(buf)

	//3.ParsePKIXPublicKey解析一个DER编码的公钥。
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("断言失败")
	}
	//使用公钥加密
	//rsa.EncryptPKCS1v15使用PKCS#1 v1.5规定的填充方案和RSA算法加密msg。
	//即使是加密完全一样的数据，产生的密文也不是相同的(但是这些密文可以解密出同样的明文数据)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	if err != nil {
		log.Fatalf("rsa.EncryptPKCS1v15()函数执行出错,错误为:%v\n", err)
	}
	return ciphertext
}

//RSA私钥解密
func RSAPrivateDecrypt(ciphertext []byte, privateKeyFile string) []byte {
	//1.打开文件并且读出文件内容
	file, err := os.Open(privateKeyFile)
	if err != nil {
		log.Fatalf("os.Open()函数执行出错,错误为:%v\n", err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("file.Stat()方法执行出错,错误为:%v\n", err)
	}
	buf := make([]byte, fileInfo.Size())
	_, err = file.Read(buf)

	//2.pem解码
	block, _ := pem.Decode(buf)

	//ParsePKCS1PrivateKey解析ASN.1 PKCS#1 DER编码的rsa私钥。
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("x509.ParsePKCS1PrivateKey()函数执行出错,错误为:%v\n", err)
	}

	//3.使用私钥解密
	//rsa.DecryptPKCS1v15使用PKCS#1 v1.5规定的填充方案和RSA算法解密密文。
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		log.Fatalf("x509.ParsePKCS1PrivateKey()函数执行出错,错误为:%v\n", err)
	}
	return plaintext
}

func main() {
	plaintext := []byte("你好吗")
	ciphertext := RSAPublicEncrypt(plaintext, "public.pem")
	plaintext = RSAPrivateDecrypt(ciphertext, "private.pem")
	fmt.Println(string(plaintext))
}
