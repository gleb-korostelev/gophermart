package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
)

func EncryptIt(value string, hash [32]byte) (string, error) {
	aesBlock, err := aes.NewCipher((hash[:]))
	if err != nil {
		return "", err
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcmInstance.NonceSize())

	return hex.EncodeToString(gcmInstance.Seal(nonce, nonce, []byte(value), nil)), nil
}

func DecryptIt(ciphered string, hash [32]byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	cephByte, _ := hex.DecodeString(ciphered)

	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := cephByte[:nonceSize], cephByte[nonceSize:]

	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return nil, err
	}

	return originalText, nil
}

func Encrypt(EncryptionKey, decryptedKey string) (string, error) {
	encryptionKey := sha256.Sum256([]byte(EncryptionKey))
	return EncryptIt(decryptedKey, encryptionKey)
}

func Decrypt(EncryptionKey, encryptedKey string) (string, error) {
	encryptionKey := sha256.Sum256([]byte(EncryptionKey))
	res, err := DecryptIt(encryptedKey, encryptionKey)
	return string(res[:]), err
}
