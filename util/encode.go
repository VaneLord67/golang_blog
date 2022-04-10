package util

import (
	"crypto/md5"
	"encoding/base64"
)

func Md5Base64Encode(data string) string {
	s := md5.Sum([]byte(data))
	sEnc := base64.StdEncoding.EncodeToString(s[:])
	return sEnc
}
