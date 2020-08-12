package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"math/rand"
)

type token struct {
	UserID string
	Rand   int
	Crc    uint32
}

func GetUserIdByToken(cookieStr string) string {
	tk := Detoken(cookieStr)
	if tk == nil {
		return ""
	}
	return tk.UserID
}

func CheckToken(cookieStr string) bool {
	tk := Detoken(cookieStr)
	if tk == nil {
		return false
	}
	crc := crc32.ChecksumIEEE([]byte(tk.UserID + fmt.Sprintf("%d", tk.Rand)))
	return crc == tk.Crc
}

func Entoken(userID string) string {
	rd := rand.Intn(10000)
	tk := token{
		UserID: userID,
		Rand:   rd,
		Crc:    crc32.ChecksumIEEE([]byte(userID + fmt.Sprintf("%d", rd))),
	}
	tkJson, err := json.Marshal(tk)
	if err != nil {
		fmt.Println(err)
	}
	return base64.StdEncoding.EncodeToString(tkJson)
}

func Detoken(cookieStr string) *token {
	tkJson, err := base64.StdEncoding.DecodeString(cookieStr)
	if err != nil {
		return nil
	}
	tk := token{}
	err = json.Unmarshal(tkJson, &tk)
	if err != nil {
		return nil
	}
	return &tk
}
