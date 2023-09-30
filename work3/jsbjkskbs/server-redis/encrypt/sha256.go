package encrypt

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(s string) string{
	sha := sha256.New()
	sha.Write([]byte(s))
	return hex.EncodeToString(sha.Sum(nil))
}