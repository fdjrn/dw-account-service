package tools

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	TransTopUp   = "Top-Up"
	TransPayment = "Payment"
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
	ciphertext, _ := hex.DecodeString(ct)
	c, err := aes.NewCipher(key)
	plain := make([]byte, len(ciphertext))

	c.Decrypt(plain, ciphertext)

	decodedStr := string(plain[:])

	result, err := strconv.Atoi(strings.TrimLeft(decodedStr, "0"))
	if err != nil {
		return 0, err
	}
	return result, err
}

func GetUnixTime() string {
	tUnixMicro := int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Microsecond)
	return strconv.FormatInt(tUnixMicro, 10)
}

func GetUnixTimeMicro() string {
	tUnixMicro := int64(time.Nanosecond) * time.Now().UnixNano() / 1000
	return strconv.FormatInt(tUnixMicro, 10)
}

func GetUnixTimeNano() string {
	tUnixMicro := int64(time.Nanosecond) * time.Now().UnixNano()
	return strconv.FormatInt(tUnixMicro, 10)
}

func GenerateReceiptNumber(transType string, id string) string {
	tUnix := GetUnixTimeNano()
	var r string

	switch transType {
	case TransTopUp:
		r = "1000" + tUnix + id
	case TransPayment:
		r = "2000" + tUnix + id
	}

	return r
}
