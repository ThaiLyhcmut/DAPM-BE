package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

// Mã hóa AES-CTR
func CreateAES(plainText string) (string, error) {
	fmt.Println(plainText)
	key := os.Getenv("AES_KEY")
	if len(key) != 32 { // AES-256 cần 32 byte
		return "", fmt.Errorf("invalid AES key size")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Tạo IV ngẫu nhiên
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// Mã hóa dữ liệu
	ciphertext := make([]byte, len(plainText))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, []byte(plainText))
	fmt.Println(hex.EncodeToString(iv) + hex.EncodeToString(ciphertext))
	// Kết hợp IV + ciphertext và encode thành hex
	return hex.EncodeToString(iv) + hex.EncodeToString(ciphertext), nil
}

func ParseASE(cipherText string) (string, error) {
	fmt.Println(cipherText)
	key := os.Getenv("AES_KEY")
	if len(key) != 32 {
		return "", fmt.Errorf("invalid AES key size")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Giải mã hex
	rawData, err := hex.DecodeString(cipherText)
	if err != nil || len(rawData) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext")
	}

	// Tách IV và dữ liệu mã hóa
	iv := rawData[:aes.BlockSize]
	ciphertext := rawData[aes.BlockSize:]

	// Giải mã dữ liệu
	plainText := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plainText, ciphertext)
	fmt.Print(string(plainText))
	return string(plainText), nil
}
