package shortener

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

const urlHashPostfix = "Xa1"

func Hash(inputString string, postfixCount int) string {
	var hash [16]byte

	hash = md5.Sum([]byte(inputString + strings.Repeat(urlHashPostfix, postfixCount)))

	return hex.EncodeToString(hash[:])
}

func Shorten(hashToShorten string) string {
	if len(hashToShorten) < 8 {
		return hashToShorten
	}

	return hashToShorten[:7]
}
