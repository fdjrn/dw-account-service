package str

import (
	"fmt"
	"github.com/dw-account-service/internal/utilities"
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandomString(length int, prefix, suffix string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = utilities.Charset[rand.Intn(len(utilities.Charset))]
	}
	return fmt.Sprintf("%s%s%s", prefix, string(b), suffix)
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

func GenerateReceiptNumber(transType int, id string) string {
	tUnix := GetUnixTimeNano()
	var r string

	switch transType {
	case utilities.TransTypeTopUp:
		r = fmt.Sprintf("1000%s%s", tUnix, id)
	case utilities.TransTypePayment:
		r = fmt.Sprintf("2000%s%s", tUnix, id)
	case utilities.TransTypeDistribution:
		r = fmt.Sprintf("3000%s%s", tUnix, id)
	}

	return r
}

func GenerateTransNumber() string {
	pre := time.Now().Format("20060102")
	return GenerateRandomString(8, pre, "")

	//rand.Seed(time.Now().UnixNano())
	//b := make([]byte, 8)
	//for i := range b {
	//	b[i] = charset[rand.Intn(len(charset))]
	//}
	//return fmt.Sprintf("%s%s", time.Now().Format("20060102"), string(b))
}
