package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
    "io"
    "os"
)

func Encrypt(plainText []byte) ([]byte, error) {
    block, err := aes.NewCipher([]byte(os.Getenv("APP_KEY")))
    if err != nil {
        return nil, err
    }

    // Generate a random IV
    cipherText := make([]byte, aes.BlockSize+len(plainText))
    iv := cipherText[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }

    // Encrypt the plaintext
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

    return cipherText, nil
}

func Decrypt(cipherText []byte) ([]byte, error) {
    block, err := aes.NewCipher([]byte(os.Getenv("APP_KEY")))
    if err != nil {
        return nil, err
    }

    if len(cipherText) < aes.BlockSize {
        return nil, fmt.Errorf("cipher text too short")
    }

    iv := cipherText[:aes.BlockSize]
    cipherText = cipherText[aes.BlockSize:]

    // Decrypt the ciphertext
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(cipherText, cipherText)

    return cipherText, nil
}
