package util

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"reflect"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func FormatPhoneNumberLocal(num string) (string, error) {
	num = strings.Replace(num, " ", "", -1)
	num = strings.Replace(num, "-", "", -1)
	num = strings.Replace(num, "+", "", -1)
	num = strings.Trim(num, "")
	if len(num) < 6 {
		return "", errors.New("Mobile number not long enough!")
	} else if strings.HasPrefix(num, "27") {
		num := strings.Replace(num, "27", "0", 1)
		// log.Printf("%s -> %s\n", num, fNum)
		return num, nil
	} else if strings.HasPrefix(num, "0") {
		return num, nil
	}
	return "", errors.New("Mobile number incorrect format!")
}

func FormatPhoneNumberToInternational(num string) (string, error) {
	num = strings.Replace(num, " ", "", -1)
	num = strings.Replace(num, "-", "", -1)
	num = strings.Trim(num, "")
	if len(num) < 6 {
		return "", errors.New("Mobile number not long enough!")
	} else if strings.HasPrefix(num, "27") {
		num := strings.Replace(num, "27", "+27", 1)
		return num, nil
	} else if strings.HasPrefix(num, "+27") {
		return num, nil
	} else if strings.HasPrefix(num, "0") {
		num := strings.Replace(num, "0", "+27", 1)
		// log.Printf("%s -> %s\n", num, fNum)
		return num, nil
	}
	return "", errors.New("Mobile number incorrect format!")
}

func GenerateCode() uint16 {
	rNo, _ := rand.Int(rand.Reader, big.NewInt(8999))
	return uint16(1000 + rNo.Int64())
}

func IndexOfItemInSlice[T any](item T, slice []T) int {
	for i, v := range slice {
		if reflect.DeepEqual(v, item) {
			return i
		}
	}
	return -1
}

func GenPasswordHash(password string) (string, error) {
	if len(strings.TrimSpace(password)) < 6 {
		return "", errors.New("Password to short")
	}
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hashedPw), err
}

func CheckPasswordMatch(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("Cannot parse empty string")
	}
	postfix := s[len(s)-1:]

	if postfix == "d" {
		days, err := time.ParseDuration(s[:len(s)-1] + "h")
		return days * 24, err
	}
	dur, err := time.ParseDuration(s)
	return dur, err
	// return time.ParseDuration(s)
}

func JsonMapFromBytes(b []byte) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m
}
