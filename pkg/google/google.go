//nolint
package google

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"strings"
	"time"
)

func hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func toBytes(value int64) []byte {
	result := []byte{}
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func oneTimePassword(key, data []byte) uint32 {
	hash := hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] &= 0x7F
	number := toUint32(hashParts)
	return number % 1000000
}

func VerifyCode(secret, code string) (bool, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := base32decode(secretUpper)
	if err != nil {
		return false, err
	}

	if code == fmt.Sprintf("%06d", oneTimePassword(secretKey, toBytes(time.Now().UTC().Unix()/30))) {
		return true, nil
	} else if code == fmt.Sprintf("%06d", oneTimePassword(secretKey, toBytes((time.Now().UTC().Unix()-30)/30))) {
		return true, nil
	} else if code == fmt.Sprintf("%06d", oneTimePassword(secretKey, toBytes((time.Now().UTC().Unix()+30)/30))) {
		return true, nil
	}

	return false, nil
}
