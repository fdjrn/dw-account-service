package tools

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	rand2 "math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	TransTopUp   = "Top-Up"
	TransPayment = "Payment"
)

var charset = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

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
		r = fmt.Sprintf("1000%s%s", tUnix, id)
	case TransPayment:
		r = fmt.Sprintf("2000%s%s", tUnix, id)
	}

	return r
}

func GenerateTransNumber() string {
	rand2.Seed(time.Now().UnixNano())
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand2.Intn(len(charset))]
	}
	return fmt.Sprintf("%s%s", time.Now().Format("20060102"), string(b))
}
