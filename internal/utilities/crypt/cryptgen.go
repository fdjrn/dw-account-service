package crypt

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func GenerateSecretKey() (string, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", key), nil
}

func Encrypt(key []byte, plaintext string) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}
	out := make([]byte, len(plaintext))
	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out), err
}

func Decrypt(key []byte, ct string) (string, error) {
	ciphertext, _ := hex.DecodeString(ct)
	c, err := aes.NewCipher(key)
	plain := make([]byte, len(ciphertext))
	c.Decrypt(plain, ciphertext)
	s := string(plain[:])

	return s, err
}

func DecryptAndConvert(key []byte, ct string) (int, error) {
	//ciphertext, _ := hex.DecodeString(ct)
	//c, err := aes.NewCipher(key)
	//plain := make([]byte, len(ciphertext))
	//
	//c.Decrypt(plain, ciphertext)
	//
	//decodedStr := string(plain[:])

	decodedStr, _ := Decrypt(key, ct)
	result, err := strconv.Atoi(strings.TrimLeft(decodedStr, "0"))
	if err != nil {
		return 0, err
	}
	return result, err
}
