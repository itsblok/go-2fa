package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	secretSize = 20
	timeStep   = 30
	codeDigits = 6
)

// GenerateSecret creates a new base32 encoded TOTP secret.
func GenerateSecret() (string, error) {
	buf := make([]byte, secretSize)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	secret := base32.StdEncoding.EncodeToString(buf)
	secret = strings.TrimRight(secret, "=")

	return secret, nil
}

// GenerateCode generates the current TOTP code.
func GenerateCode(secret string, t time.Time) (string, error) {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("invalid secret encoding")
	}

	counter := uint64(t.Unix() / timeStep)

	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], counter)

	h := hmac.New(sha1.New, key)
	h.Write(msg[:])

	hash := h.Sum(nil)
	offset := hash[len(hash)-1] & 0x0f

	code := (int(hash[offset])&0x7f)<<24 |
		int(hash[offset+1])<<16 |
		int(hash[offset+2])<<8 |
		int(hash[offset+3])

	code %= 1000000

	return fmt.Sprintf("%06d", code), nil
}

// VerifyCode validates a TOTP code.
func VerifyCode(secret, input string, t time.Time) bool {
	code, err := GenerateCode(secret, t)
	if err != nil {
		return false
	}

	return code == input
}

// BuildOTPAuthURL builds an otpauth URL for authenticator apps.
func BuildOTPAuthURL(issuer, account, secret string) string {
	label := url.QueryEscape(issuer + ":" + account)

	return fmt.Sprintf(
		"otpauth://totp/%s?secret=%s&issuer=%s",
		label,
		secret,
		url.QueryEscape(issuer),
	)
}
