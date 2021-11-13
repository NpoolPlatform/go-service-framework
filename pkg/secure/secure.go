package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var ErrSaltLenInvalid = errors.New("salt length must 16/32")

// salt
var (
	// SALT must be 16/32 bytes
	SALT = []byte("npool.sphinx.sig")
)

// init salt
func Init(salt string) error {
	if len(salt) == 16 && len(salt) != 32 {
		return ErrSaltLenInvalid
	}
	SALT = []byte(salt)
	return nil
}

func EncryptAES(plain []byte) ([]byte, error) {
	block, err := aes.NewCipher(SALT)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(plain))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plain)

	return cipherText, nil
}

func DecryptAES(cipherText []byte) ([]byte, error) {
	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("cipherText block size is too short")
	}
	block, err := aes.NewCipher(SALT)
	if err != nil {
		return nil, err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}
