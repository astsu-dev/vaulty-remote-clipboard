package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

const pbkdf2Salt = "super-secret-vaulty-salt"
const pbkdf2Iterations = 600000
const pbkdf2KeyLen = 32

// Derives pbkdf2 key from the specified password
func DerivePbkdf2From(password []byte) []byte {
	return pbkdf2.Key(
		password,
		[]byte(pbkdf2Salt),
		pbkdf2Iterations,
		pbkdf2KeyLen,
		sha256.New,
	)
}

// DecryptGCM decrypts ciphertext in base64 encoding with the key.
// If message authentication fails will return error.
func DecryptGCM(key []byte, ciphertext []byte) ([]byte, error) {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(ciphertext)))
	n, err := base64.StdEncoding.Decode(decoded, ciphertext)
	if err != nil {
		panic(err)
	}
	ciphertext = decoded[:n]

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Encrypts text with the key. Returns result in the base64 encoding.
func EncryptGCM(key []byte, text []byte) []byte {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	encryptedData := gcm.Seal(nonce, nonce, text, nil)

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(encryptedData)))
	base64.StdEncoding.Encode(encoded, encryptedData)

	return encoded
}
