package sherbet

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)


func Sign(input string, key string, salt string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(fmt.Sprintf("%s|%s", salt, input)))

	result := mac.Sum(nil)

	return hex.EncodeToString(result)
}

func VerfiySignature(input string, inputMAC string, key string, salt string) bool {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(fmt.Sprintf("%s|%s", salt, input)))
	expectedMAC := mac.Sum(nil)
	decodedHash, _ := hex.DecodeString(inputMAC)
	return hmac.Equal(decodedHash, expectedMAC)
}
