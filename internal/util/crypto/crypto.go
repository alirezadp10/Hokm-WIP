package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// Encrypt encrypts the plaintext and returns a Base64-encoded ciphertext
func Encrypt(plainText string) (string, error) {
	// Ensure the app key is of the correct length
	appKey := os.Getenv("APP_KEY")
	if len(appKey) != 32 {
		return "", errors.New("APP_KEY must be 32 bytes long for AES-256")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(appKey))
	if err != nil {
		return "", err
	}

	// Generate a random IV
	cipherText := make([]byte, aes.BlockSize+len([]byte(plainText)))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt the plaintext
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))

	// Encode the ciphertext as Base64
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts a Base64-encoded ciphertext and returns the plaintext
func Decrypt(cipherTextBase64 string) (string, error) {
	// Ensure the app key is of the correct length
	appKey := os.Getenv("APP_KEY")
	if len(appKey) != 32 {
		return "", errors.New("APP_KEY must be 32 bytes long for AES-256")
	}

	// Decode the Base64-encoded ciphertext
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(appKey))
	if err != nil {
		return "", err
	}

	// Extract the IV from the ciphertext
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	// Decrypt the ciphertext
	stream := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
