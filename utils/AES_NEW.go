package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const(
	sKey        = "dde4b1f8a9e6b814"
	ivParameter     = "test121412412312"
)


func PswDecrypt(src string)(string)  {
	key := []byte(sKey)
	iv := []byte(ivParameter)

	var result []byte
	var err error

	result,err=base64.RawStdEncoding.DecodeString(src)
	if err != nil {
		panic(err)
	}
	origData, err := Aes128Decrypt(result, key, iv)
	if err != nil {
		panic(err)
	}
	return string(origData)
}

func Aes128Decrypt(crypted []byte,key []byte,IV []byte)([]byte,error){
	block,err := aes.NewCipher(key)
	if err != nil{
		return nil,err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block,IV[:blockSize])

	origData := make([]byte,len(crypted))
	blockMode.CryptBlocks(origData,crypted)
	origData = PKCS5UnPadding(origData)
	return origData,nil
}

/*func PKCS5UnPadding(origData []byte) []byte  {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}*/
