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

	var input string
	fmt.Print("Enter code: ")
	fmt.Scanln(&input)

	if totp.VerifyCode(secret, input, time.Now()) {
		fmt.Println("Valid code")
	} else {
		fmt.Println("Invalid code")
	}
}
