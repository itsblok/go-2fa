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
	secretSize = 20 // 160-bit secret
	timeStep   = 30 // seconds
	digits     = 6
)

// GenerateSecret returns a new Base32-encoded secret.
func GenerateSecret() (string, error) {
	b := make([]byte, secretSize)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// Base32 without padding for compatibility with authenticator apps
	return strings.TrimRight(base32.StdEncoding.EncodeToString(b), "="), nil
}

// GenerateCode produces a TOTP code for a given secret and timestamp.
func GenerateCode(secret string, t time.Time) (string, error) {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}

	counter := uint64(t.Unix() / timeStep)

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], counter)

	h := hmac.New(sha1.New, key)
	h.Write(buf[:])
	hash := h.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F

	truncated := binary.BigEndian.Uint32(hash[offset : offset+4])
	truncated &= 0x7FFFFFFF

	code := truncated % 1000000

	return fmt.Sprintf("%06d", code), nil
}

// BuildOTPAuthURL creates an otpauth URL compatible with authenticator apps.
func BuildOTPAuthURL(issuer, account, secret string) string {
	label := url.QueryEscape(fmt.Sprintf("%s:%s", issuer, account))
	issuerParam := url.QueryEscape(issuer)

	return fmt.Sprintf(
		"otpauth://totp/%s?secret=%s&issuer=%s&period=30&digits=6",
		label,
		secret,
		issuerParam,
	)
}

// VerifyCode checks if the provided code is valid within a time window.
func VerifyCode(secret, code string, t time.Time) bool {
	// Allow ±1 time step to tolerate clock drift
	for i := -1; i <= 1; i++ {
		ts := t.Add(time.Duration(i*timeStep) * time.Second)
		gen, err := GenerateCode(secret, ts)
		if err != nil {
			return false
		}
		if gen == code {
			return true
		}
	}
	return false
}
