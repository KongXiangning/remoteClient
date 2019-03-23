package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

//var key = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi11c2VyLXRva2VuLXpzbTY4Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImFkbWluLXVzZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJlOGQ3ZWU5ZC1lM2ZmLTExZTgtYTQwNS1mYTE2M2VmN2RjZTEiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06YWRtaW4tdXNlciJ9.nHaojEVO2uRyB9tPR61HRnA2kcOYLyOohIxESohWmD9NvbBsitNMq9UTZ3Avwp0lI__HAmK7VyByCJiF3Ornhhw2sTJAKvs8kgWlPwy0WaJtsQQRFvw7uYIN0mKMVA0iaUGXkiJXMk-z8zktl5t6vvdhj10ug4pd2-zv1vI7YbvAAIQTop2w3yR7xz7XWEUD5K8AIuwekx9X6sPu3EztlDHBfg4OsPAGc7VO21huH_Zn5qCjDF4AIwi2XvvvS7ZIckiODy2YzTWStodJL6pdyohnSVjh2rEmEP60XOH3sZZGoiCFRDETRWPTBRORrkkq2gydBdSxGerVznQPAiP8yA"
var key = "dde4b1f8a9e6b814"
var iv = "test121412412312"

func Encrypt(origData []byte)([]byte,error)  {

	block, err := aes.NewCipher([]byte(key))
	if err != nil{
		return nil,err
	}

	blockSize := block.BlockSize()
	origData = pkcs5Padding(origData,blockSize)

	blockMode := cipher.NewCBCEncrypter(block,[]byte(iv))
	crypted := make([]byte,len(origData))

	blockMode.CryptBlocks(crypted,origData)
	return crypted,err
}

func Decrypt(crypted []byte)([]byte,error)  {
	block,err := aes.NewCipher([]byte(key))
	if err != nil{
		fmt.Println(err)
		return nil,err
	}

	blockMode := cipher.NewCBCDecrypter(block,[]byte(iv))
	origData := make([]byte,len(crypted))
	blockMode.CryptBlocks(origData,crypted)

	origData = PKCS5UnPadding(origData)
	return origData,err

}

func pkcs5Padding(ciphertext []byte,blockSize int) []byte  {
	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)},padding)
	return append(ciphertext,padtext...)
}

func PKCS5UnPadding(origData []byte) []byte  {
	length := len(origData)
	unpadding := int(origData[length - 1])
	return origData[:(length - unpadding)]
}
