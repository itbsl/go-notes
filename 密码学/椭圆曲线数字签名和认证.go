package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
)

//使用椭圆曲线进行数字签名和认证
//1.生成秘钥对，并将秘钥对(公钥、私钥)保存到磁盘中
//2.使用私钥进行数字签名
//3.使用公钥验证数字签名

func GenerateEccKey() {
	//================== 生成私钥 ==================
	//1.使用ecdsa生成秘钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("ecdsa.GenerateKey()函数执行出错,错误为:%v\n", err)
	}

	//2.将私钥写入磁盘
	//使用x509进行序列化
	derText, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("x509.MarshalECPrivateKey()函数执行出错,错误为:%v\n", err)
	}
	//将得到的切片字符串放入到pem.Block结构体中
	block := pem.Block{
		Type:    "ECDSA PRIVATE KEY",
		Headers: nil,
		Bytes:   derText,
	}
	//使用pem编码
	privateFile, err := os.Create("ecc_private.pem")
	if err != nil {
		log.Fatalf("os.Create()函数执行出错,错误为:%v\n", err)
	}
	err = pem.Encode(privateFile, &block)
	if err != nil {
		log.Fatalf("pem.Encode()函数执行出错,错误为:%v\n", err)
	}
	privateFile.Close()

	//================== 生成公钥 ==================
	//3.将公钥写入磁盘
	//从私钥中得到公钥
	publicKey := privateKey.PublicKey
	//使用x509进行序列化
	derText, err = x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		log.Fatalf("x509.MarshalPKIXPublicKey()函数执行出错,错误为:%v\n", err)
	}

	//将得到的切片字符串放入到pem.Block结构体中
	block = pem.Block{
		Type:    "ECDSA PUBLIC KEY",
		Headers: nil,
		Bytes:   derText,
	}

	//使用pem进行编码
	publicFile, err := os.Create("ecc_public.pem")
	if err != nil {
		log.Fatalf("os.Create()函数执行出错,错误为:%v\n", err)
	}
	err = pem.Encode(publicFile, &block)
	if err != nil {
		log.Fatalf("pem.Encode()函数执行出错,错误为:%v\n", err)
	}
	publicFile.Close()
}

//签名
func EccSign(plaintext []byte, privateKeyFile string) (rText, sText []byte) {
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

	//4.使用x509对私钥进行还原
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("x509.ParseECPrivateKey()函数执行出错,错误为:%v\n", err)
	}

	//5.对原始数据进行hash运算 -> 散列值
	hashText := sha256.Sum256(plaintext)

	//6.进行数字签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashText[:])
	if err != nil {
		log.Fatalf("ecdsa.Sign()函数执行出错,错误为:%v\n", err)
	}

	//7.对r, s中的数据进行格式化，转换成[]byte
	rText, err = r.MarshalText()
	if err != nil {
		log.Fatalf("r.MarshalText()方法执行出错,错误为:%v\n", err)
	}
	sText, err = s.MarshalText()
	if err != nil {
		log.Fatalf("s.MarshalText()方法执行出错,错误为:%v\n", err)
	}
	return rText, sText
}

//验签
func VerifyECCSignature(plaintext, rText, sText []byte, publicKeyFile string) bool {
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

	//3.使用pem解码 -> 得到pem.Block结构体变量
	block, _ := pem.Decode(buf)

	//4.使用x509对公钥还原
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("x509.ParsePKIXPublicKey()函数执行出错,错误为:%v\n", err)
	}
	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalln("断言失败")
	}

	//5.对原始数据进行哈希运算 -> 得到散列值
	hashText := sha256.Sum256(plaintext)

	//6.将rTest, sText再转换成*big.Int型数据
	var r, s big.Int
	r.UnmarshalText(rText)
	s.UnmarshalText(sText)

	//7.认证签名
	return ecdsa.Verify(publicKey, hashText[:], &r, &s)
}

func main() {
	//1.生成秘钥对，并将秘钥对(公钥、私钥)保存到磁盘中
	GenerateEccKey()

	//2.使用私钥进行数字签名
	plaintext := []byte("hello world")
	rText, sText := EccSign(plaintext, "ecc_private.pem")

	//3.使用公钥验证数字签名
	fmt.Println(VerifyECCSignature(plaintext, rText, sText, "ecc_public.pem"))
}