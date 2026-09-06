package base62

import (
	"errors"
	"strings"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base = uint64(len(alphabet)) // 62

func Encode(n uint64) string {
	if n == 0 {
		return string(alphabet[0])
	}

	var sb strings.Builder
	for n > 0 {
		remainder := n % base
		sb.WriteByte(alphabet[remainder])
		n /= base
	}

	encoded := []byte(sb.String())
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}
	return string(encoded)
}

func Decode(s string) (uint64, error) {
	var n uint64
	for _, c := range s {
		idx := strings.IndexRune(alphabet, c)
		if idx < 0 {
			return 0, errors.New("base62: invalid character in input")
		}
		n = n*base + uint64(idx)
	}
	return n, nil
}
