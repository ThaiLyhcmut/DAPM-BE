package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"os"
)

func CreateASE(plainText string) string {
	key := os.Getenv("AES_KEY")
	block, _ := aes.NewCipher([]byte(key))
	ciphertext := make([]byte, len(plainText))
	stream := cipher.NewCTR(block, []byte(key[:aes.BlockSize]))
	stream.XORKeyStream(ciphertext, []byte(plainText))
	return hex.EncodeToString(ciphertext)
}

func ParseASE(cipherText string) string {
	key := os.Getenv("AES_KEY")
	block, _ := aes.NewCipher([]byte(key))
	ciphertext, _ := hex.DecodeString(cipherText)
	plainText := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, []byte(key[:aes.BlockSize]))
	stream.XORKeyStream(plainText, ciphertext)
	return string(plainText)
}
