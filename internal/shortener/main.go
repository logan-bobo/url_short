package shortener

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

const urlHashPostfix = "Xa1"

func hash(inputString string, withPostfix bool, postfixCount int) string {
	var hash [16]byte

	if withPostfix {
		hash = md5.Sum([]byte(inputString + strings.Repeat(urlHashPostfix, postfixCount)))
	} else {
		hash = md5.Sum([]byte(inputString))
	}

	return hex.EncodeToString(hash[:])
}

func shorten(hashToShorten string) string {
	if len(hashToShorten) < 8 {
		return hashToShorten
	}

	return hashToShorten[:7]
}
