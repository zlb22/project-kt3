package kit

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/forgoer/openssl"
)

var secretError = errors.New("secret key lost or too short")
var dataError = errors.New("not encrypted data")

func Decrypt(data, key string) ([]byte, error) {
	if len(key) < 16 {
		return nil, secretError
	}
	keyByte := []byte(key)
	iv := keyByte[:16]
	res, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	if len(res)%16 != 0 {
		return nil, dataError
	}
	return openssl.AesCBCDecrypt([]byte(res), keyByte, iv, openssl.PKCS7_PADDING)
}

func Encrypt(data interface{}, key string) (string, error) {
	if len(key) < 16 {
		return "", secretError
	}
	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	keyByte := []byte(key)
	iv := keyByte[:16]
	dst, err := openssl.AesCBCEncrypt(res, keyByte, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}
	result := base64.StdEncoding.EncodeToString(dst)
	return result, err
}

//DecryptStr 解密字符串，和EncryptStr配套使用
func DecryptStr(data, key string) (string, error) {
	if len(key) < 16 {
		return "", secretError
	}
	keyByte := []byte(key)
	iv := keyByte[:16]
	res, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	if len(res)%16 != 0 {
		return "", dataError
	}
	strByte, err := openssl.AesCBCDecrypt([]byte(res), keyByte, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}
	return string(strByte), nil
}

// EncryptStr 加密字符串，和DecryptStr配套使用
func EncryptStr(data, key string) (string, error) {
	if len(key) < 16 {
		return "", secretError
	}
	dataByte := []byte(data)
	keyByte := []byte(key)
	iv := keyByte[:16]
	dst, err := openssl.AesCBCEncrypt(dataByte, keyByte, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}
	result := base64.StdEncoding.EncodeToString(dst)
	return result, err
}
