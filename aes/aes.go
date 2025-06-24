package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"strconv"
)

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}
func newECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func pKCS7UnPadding(data []byte, blockSize int) ([]byte, error) {
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, errors.New("invalid padding")
	}

	padding := int(data[len(data)-1])
	if padding > blockSize || padding == 0 {
		return nil, errors.New("invalid padding")
	}

	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:len(data)-padding], nil
}
func ECBDecrypt(ciphertext string, key []byte) (string, error) {
	decodeString, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockMode := newECBDecrypter(block)
	word := make([]byte, len(decodeString))
	blockMode.CryptBlocks(word, decodeString)
	unpaddedPlaintext, err := pKCS7UnPadding(word, block.BlockSize())
	if err != nil {
		return "", err
	}
	return string(unpaddedPlaintext), nil
}

func GenerateSaltedHash(password string, salt int64) string {
	// 将密码和盐进行拼接
	str := strconv.FormatInt(salt, 10)
	saltedPassword := password + str
	// 创建一个MD5哈希对象
	hash := md5.New()
	// 将拼接后的字符串转换为字节数组并计算哈希值
	_, _ = io.WriteString(hash, saltedPassword)
	hashedPassword := hash.Sum(nil)
	// 将哈希值转换为十六进制字符串表示
	hashedPasswordHex := hex.EncodeToString(hashedPassword)
	return hashedPasswordHex
}
