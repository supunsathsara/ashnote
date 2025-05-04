package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt encrypts the plaintext with the provided password
func Encrypt(plaintext, password string) (string, error) {
	// Create a new hash from the password
	key := sha256.Sum256([]byte(password))

	// Create a new AES cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure
	// So we include it at the beginning of the ciphertext
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Use CFB mode for encryption
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	// Return base64 encoded string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the ciphertext with the provided password
func Decrypt(encryptedText, password string) (string, error) {
	// Decode the base64 encoded string
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// Create a new hash from the password
	key := sha256.Sum256([]byte(password))

	// Create a new AES cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// Check if the ciphertext is valid
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Get the IV from the beginning of the ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Use CFB mode for decryption
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
