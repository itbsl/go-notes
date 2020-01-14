package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

//CTR - CounTeR,计数器模式
//1.特点:密文没有规律，明文分组是和一个数据流进行的按位异或操作，最终生成了密文
//2.不需要初始化向量
//	go接口中的iv可以理解为随机数种子，iv的长度==明文分组的长度
//3.不需要填充


//AES加密
//plainText 明文
//key 秘钥，大小为16byte
func AESCTREncrypt(plaintext []byte, key []byte) (ciphertext []byte, err error) {
	//1.创建并返回一个使用DES算法的cipher.Block接口
	//	秘钥长度为128bit,即128/8 = 16字节(byte)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//2.创建一个计数器模式的、底层采用block生成key流的Stream接口，
	//初始向量iv的长度必须等于block的块尺寸。
	//vi向量必须是独一无二的，但是不需要是加密的
	//(加密本是不需要向量的，但是加密需要随机数和计数器，随机数是通过把向量作为种子生成的随机数)
	ciphertext = make([]byte, aes.BlockSize + len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)

	//3.从加密器的key流和src中依次取出字节二者xor后写入dst，src和dst可指向同一内存地址
	//最终返回的实际加密结果ciphertext中,前blockSize()个字节其实是iv向量，后面才是真正加密的结果
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	//5.返回结果
	return
}

//cipherText 密文
//key 秘钥，和加密秘钥相同,大小为:16byte
func AESCTRDecrypt(ciphertext, key []byte) (plaintext []byte, err error) {

	//1.创建并返回一个使用AES算法的cipher.Block接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	///2.创建一个计数器模式的、底层采用block生成key流的Stream接口，
	//初始向量iv的长度必须等于block的块尺寸。
	//vi向量必须是独一无二的，但是不需要是加密的
	//创建明文切片，切片长度为密文减去blockSize()长度即可，因为密文中有blockSize()长度的原始向量iv
	plaintext = make([]byte, len(ciphertext) - block.BlockSize())
	iv := ciphertext[:aes.BlockSize] //取出原始向量
	stream := cipher.NewCTR(block, iv)

	//3.解密数据(只需要解密ciphertext中blockSize()后的数据，因为前面blockSize()个长度是iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	//4.返回结果
	return
}

func main() {

	//明文
	plaintext := []byte("你好呀")
	//秘钥(长度必须为16字节，因为AES加密就是按照16字节(128bit)一组加密的，规定秘钥长度和每组要加密的长度相等)
	key := []byte("1234567890123456")
	ciphertext, err := AESCTREncrypt(plaintext, key)
	if err != nil {
		log.Fatalf("加密出错,错误为:%v\n", err)
	}
	fmt.Println(string(ciphertext))

	cipherDecryptText, err := AESCTRDecrypt(ciphertext, key)
	if err != nil {
		log.Fatalf("解密出错,错误为:%v\n", err)
	}
	fmt.Println(string(cipherDecryptText))
}