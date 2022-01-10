package cryp

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"strings"
)

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

// AESCBCD  AES-CBC算法解密函数
//
//参数为密文，密钥，iv
func AESCBCD(encryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	// 解填充
	encryptData, _ = pkcs7Unpad(encryptData, block.BlockSize())
	return encryptData, err
}

// AECCRTE AES-CRT算法加密函数
//
// 将需要加密的数据、加密密钥key、加密向量vi作为参数进行加密，函数内置了pkcs7 填充数据填充模式，返回数据为加密后的数据(字符串类型十六进制)
func AECCRTE(data []byte, key []byte, vi []byte) (string, error) {

	//指定加密、解密算法为AES，返回一个AES的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	text, err := pkcs7Pad(data, block.BlockSize()) // 补全码
	if err != nil {
		return "", err
	}

	//指定分组模式
	blockMode := cipher.NewCTR(block, vi[:16])

	//执行加密、解密操作
	message := make([]byte, len(text))
	blockMode.XORKeyStream(message, text)

	//将byte转为其16进制的字符串
	encodedStr := hex.EncodeToString(message)
	return strings.ToUpper(encodedStr), nil
}

// AESCRTD AES-CRT算法解密函数
//
// 需要将 需要解密的数据(字符串类型的十六进制)、加密密钥key、加密向量vi 传递给函数，函数使用AES-CRT解密并使用pkcs7模式进行去掉补全码，返回解密后的数据
func AESCRTD(text string, key []byte, vi []byte) ([]byte, error) {

	//16进制的字符串转为对应的byte
	aa, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//指定分组模式
	blockMode := cipher.NewCTR(block, vi[:16])

	//执行加密、解密操作
	message := make([]byte, len(aa))
	blockMode.XORKeyStream(message, aa)

	//去除补全码
	decrypted, err := pkcs7Unpad(message, block.BlockSize())
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// pkcs7Pad aes算法pkcs7方式为数据添加补全码
//
// 需要将 需要处理的数据以及数据块大小作为参数传递给函数，函数根据参数进行填充补全码，并返回最终数据
func pkcs7Pad(data []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if data == nil || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(data) % blocksize)
	pb := make([]byte, len(data)+n)
	copy(pb, data)
	copy(pb[len(data):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7Unpad Aes算法pkcs7方式验证和去除加载给定字节片中的补全码
//
// 需要将需要处理的数据以及块大小传递给函数，函数使用pkcs7模式去除补全码
func pkcs7Unpad(data []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if data == nil || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(data)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}

// AESCRTD AES-CRT算法解密函数
//
// 需要将 需要解密的数据(字符串类型的十六进制)、加密密钥key、加密向量vi 传递给函数，函数使用AES-CRT解密并使用pkcs7模式进行去掉补全码，返回解密后的数据
func AESCRTDNoPad(text string, key []byte, vi []byte) ([]byte, error) {

	//16进制的字符串转为对应的byte
	aa, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//指定分组模式
	blockMode := cipher.NewCTR(block, vi[:16])

	//执行加密、解密操作
	message := make([]byte, len(aa))
	blockMode.XORKeyStream(message, aa)

	return message, nil
}

// AECCRTE AES-CRT算法加密函数
//
// 将需要加密的数据、加密密钥key、加密向量vi作为参数进行加密，函数内置了pkcs7 填充数据填充模式，返回数据为加密后的数据(字符串类型十六进制)
func AECCRTENopad(data []byte, key []byte, vi []byte) (string, error) {

	//指定加密、解密算法为AES，返回一个AES的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//指定分组模式
	blockMode := cipher.NewCTR(block, vi[:16])

	//执行加密、解密操作
	message := make([]byte, len(data))
	blockMode.XORKeyStream(message, data)

	//将byte转为其16进制的字符串
	encodedStr := hex.EncodeToString(message)
	return strings.ToUpper(encodedStr), nil
}
