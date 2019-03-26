package utils

import (
	"crypto/aes"
	"crypto/cipher"
)

var SKey []byte

func Aes128Encrypt(origData []byte, key []byte, IV []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, IV[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func Aes128Decrypt(crypted []byte, key []byte, IV []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, IV[:blockSize])

	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

/*func PKCS5UnPadding(origData []byte) []byte  {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}*/
