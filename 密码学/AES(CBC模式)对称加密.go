package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
)

//CBC - Cipher Block Chaining,密码块链模式
//特点:
//1.密文没有规律，经常使用的加密方式
//2.最后一个明文分组需要填充
//	aes 最后一个分组填充满16字节
//3.需要一个初始化向量-一个数组
//	数组的长度: 与明文分组相等
//	数据来源:负责加密的人提供的
//	加解密使用的初始化向量值必须相同

//AES加密
//plainText 明文
//key 秘钥，大小为16byte
func AESCBCEncrypt(plainText []byte, key []byte) (cipherText []byte, err error) {
	//1.创建并返回一个使用DES算法的cipher.Block接口
	//	秘钥长度为128bit,即128/8 = 16字节(byte)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//2.对最后一个明文分组进行数据填充
	//*DES是以128比特的明文(比特序列)为一个单位来进行加密的
	//*最后一个如果不够128bit,则需要进行数据填充
	plainText = PKCS5Padding(plainText, block.BlockSize())

	//3.创建一个密码分组为链接模式的，底层使用DES加密的BlockMode接口
	//参数iv(向量)的长度，必须等于BlockSize
	iv := []byte("IamIVIamIVIamIVI")
	blockMode := cipher.NewCBCEncrypter(block, iv)

	//4.加密连续的数据块
	//CryptBlocks()方法指出src和dst可以指向同一个地址，所以不需要make，少分配一次内存
	//cipherText = make([]byte, len(plainText))
	cipherText = plainText
	blockMode.CryptBlocks(cipherText, plainText)

	//5.返回结果
	return
}

//使用pks5的方式填充
func PKCS5Padding(plainText []byte, blockSize int) []byte {
	//1.计算最后一个分组缺多少个字节
	//(即使能正好分组，也要填充，而且是填充一个整组的长度，
	//因为只有无论是否正好分组都填充这样方便解密的时候去除填充，否则解密的时候不知道是否填充了)
	padding := blockSize - (len(plainText) % blockSize)

	//2.创建一个大小为padding的切片，每个字节的值为padding
	//需要给plainText补足字节数，这里采用的重复数字padding padding次
	//假设padding=5,那么就相当于拼接了5个5，padText="55555"
	//这样填充还有一个好处就是解密需要去掉填充的时候，读取最后一个字节的值就知道要去掉多少个了
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	//3.将padText添加到原始数据的后边, 将最后一个分组缺少的字节数补齐
	plainText = append(plainText, padText...)

	return plainText
}

//cipherText 密文
//key 秘钥，和加密秘钥相同,大小为:8byte
func AESCBCDecrypt(cipherText, key []byte) (plainText []byte, err error) {

	//1.创建并返回一个使用AES算法的cipher.Block接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//2.创建一个密码分组为链接模式的，底层使用DES解密的BlockMode接口
	iv := []byte("IamIVIamIVIamIVI")
	blockMode := cipher.NewCBCDecrypter(block, iv)

	//3.解密数据
	//CryptBlocks()方法指出src和dst可以指向同一个地址，所以不需要make，少分配一次内存
	//plainText = make([]byte, len(cipherText))
	plainText = cipherText
	blockMode.CryptBlocks(plainText, cipherText)
	//4.去掉最后一组填充的数据
	plainText = PKCS5UnPadding(plainText)

	//5.返回结果
	return
}

//删除pks5填充的尾部数据
func PKCS5UnPadding(originData []byte) []byte {

	//1.计算数据的总长度
	length := len(originData)

	//2.根据填充的字节值得到填充的次数
	number := int(originData[length-1])

	//3.将尾部填充的number个字节去掉
	return originData[:(length - number)]
}

func main() {

	//明文
	plainText := []byte("itbsl")
	//秘钥(长度必须为16字节，因为AES加密就是按照16字节(128bit)一组加密的，规定秘钥长度和每组要加密的长度相等)
	key := []byte("1234567890123456")
	cipherText, err := AESCBCEncrypt(plainText, key)
	if err != nil {
		log.Fatalf("加密出错,错误为:%v\n", err)
	}
	fmt.Println(string(cipherText))
	cipherDecryptText, err := AESCBCDecrypt(cipherText, key)
	if err != nil {
		log.Fatalf("解密出错,错误为:%v\n", err)
	}
	fmt.Println(string(cipherDecryptText))
}