package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cheemney/go-2fa/totp"
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
	otpURL := "otpauth://totp/?secret=" + secret

	qr, err := totp.GenerateQRCode(otpURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nScan this QR with your authenticator app:")

	// terminal print (debug)
	qr.Print()

	// png export (experimental)
	if err := qr.SavePNG("qr.png", 10); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nQR saved as qr.png")

	var input string
	fmt.Print("\nEnter code: ")

	if _, err := fmt.Scanln(&input); err != nil {
		log.Fatal(err)
	}

	fmt.Println("QR Size:", qr.Size())
	fmt.Println("Mask Used:", qr.Mask())

	if totp.VerifyCode(secret, input, time.Now()) {
		fmt.Println("Valid code")
	} else {
		fmt.Println("Invalid code")
	}
}
