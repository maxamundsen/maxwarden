package main

import (
	"maxwarden/entries"
	"maxwarden/security"
	"os"

	"github.com/google/uuid"
)

func main() {
	if len(os.Args) == 2 {
		passHash, _ := security.HashPassword(os.Args[1])

		testData := []entries.Secret{}

		for range 10 {
			dummyData := entries.Secret{ID: uuid.New().String(), Description: "Twitter / X.com", URL: "https://x.com", Notes: "2fa is enabled for this account.", Username: "@johntwitter", Password: "##CORRECT_HORSE_BATTERY_STAPLE_51"}
			testData = append(testData, dummyData)
		}

		masterKey := security.SHA512_58(os.Args[1])
		cryptData, _ := security.EncryptDataWithKey(&testData, masterKey)

		println(passHash)
		println(cryptData)
	} else {
		println("Please input a password as first program argument")
	}
}
