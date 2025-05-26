package common

import (
	"encoding/base64"

	"github.com/google/uuid"
)

func Encryptbase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func DecryptBase64(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func GenerateUUID() string {
	return uuid.New().String()
}
