package libs

import (
	"crypto/aes"
	"crypto/cipher"
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func AesEncrypt(data []byte) ([]byte, error) {
	//fmt.Println(commonIV)

	//aes的加密字符串
	key_text := "astaxie12798akljzmknm.ahkjkljl;k"

	// 创建加密算法aes
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		return nil, err
	}

	//加密字符串
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(data))
	cfb.XORKeyStream(ciphertext, data)

	//fmt.Printf("%s=>%x\n", data, ciphertext)
	return ciphertext, nil
}

func AesDencrypt(data []byte) ([]byte, error) {
	//aes的加密字符串
	key_text := "astaxie12798akljzmknm.ahkjkljl;k"

	// 创建加密算法aes
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		return nil, err
	}

	// 解密字符串
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(data))
	cfbdec.XORKeyStream(plaintextCopy, data)
	//fmt.Printf("%x=>%s\n", data, plaintextCopy)

	return plaintextCopy, nil
}
