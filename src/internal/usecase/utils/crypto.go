package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
)

func GenerateSecret() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func Encrypt(secret string, plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encrypt := aesGCM.Seal(nil, nonce, []byte(plainText), nil)
	return base64.URLEncoding.EncodeToString(append(nonce, encrypt...)), nil
}

func Decrypt(secret string, encryptedStr string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("invalid encrypted text")
	}

	nonce, plaintext := data[:nonceSize], data[nonceSize:]
	decrypted, err := aesGCM.Open(nil, nonce, plaintext, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func GenerateSecureCode(min, max int) (int, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}
	return int(num.Int64()) + min, nil
}
