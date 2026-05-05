package main

import (
	"fmt"
	"log"
	"time"

	"simple-go-2fa/totp"
)

func main() {
	secret, err := totp.GenerateSecret()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Secret:", secret)

	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Current Code:", code)

	// Build OTP URL (used by authenticator apps)
	otpURL := totp.BuildOTPAuthURL("SimpleGo2FA", "user@example.com", secret)

	// Generate QR (terminal output)
	qr, err := totp.GenerateQRCode(otpURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nScan this QR with your authenticator app:")
	qr.Print()

	var input string
	fmt.Print("\nEnter code: ")
	fmt.Scanln(&input)

	if totp.VerifyCode(secret, input, time.Now()) {
		fmt.Println("Valid code")
	} else {
		fmt.Println("Invalid code")
	}
}
