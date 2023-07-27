package verify

import (
	"encoding/base64"
	"financial_statement/internal/apiserver/dal/model"

	"golang.org/x/crypto/argon2"
)

func VerifyPassword(user *model.User, password string) bool {
	return CalcPassword(password, user.Salt) == user.Password
}

func CalcPassword(password, salt string) string {
	return string(base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)))
}
