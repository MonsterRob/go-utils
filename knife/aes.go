package knife

// Only ECB/CBC aes encryption/decryption supported
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type AesOp string

const (
	CBC AesOp = "CBC"
	ECB AesOp = "ECB"
)

type Coding int

const (
	_ Coding = iota
	HEX
	BASE64
)

func AesEncrypt(plain, key []byte, coding Coding, op AesOp) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	var res []byte
	switch op {
	case ECB:
		size := block.BlockSize()
		plain = PKCS7Padding(plain, size)
		res = make([]byte, len(plain))
		for bs, be := 0, size; bs < len(plain); bs, be = bs+size, be+size {
			block.Encrypt(res[bs:be], plain[bs:be])
		}
	case CBC:
		size := aes.BlockSize
		plain = PKCS7Padding(plain, size)
		res = make([]byte, aes.BlockSize+len(plain))
		iv := res[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			panic(err)
		}
		fmt.Println(string(iv))
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(res[aes.BlockSize:], plain)
	default:
		return nil, errors.New("op invalid")
	}
	switch coding {
	case HEX:
		res = []byte(hex.EncodeToString(res))
	case BASE64:
		res = []byte(base64.StdEncoding.EncodeToString(res))
	}

	return res, nil
}

func AesDecrypt(data, key []byte, coding Coding, op AesOp) ([]byte, error) {
	var err error
	var res []byte
	switch coding {
	case HEX:
		data, err = hex.DecodeString(string(data))
	case BASE64:
		data, err = base64.StdEncoding.DecodeString(string(data))
	}
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	switch op {
	case ECB:
		size := block.BlockSize()
		res = make([]byte, len(data))
		for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
			block.Decrypt(res[bs:be], data[bs:be])
		}
		res = PKCS7UnPadding(res)
	case CBC:
		if len(data) < aes.BlockSize {
			panic("data too short")
		}
		iv := data[:aes.BlockSize]
		mode := cipher.NewCBCDecrypter(block, iv)
		data = data[aes.BlockSize:]
		res = make([]byte, len(data))
		mode.CryptBlocks(res, data)
		res = PKCS7UnPadding(res)
	}

	return res, nil

}
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
