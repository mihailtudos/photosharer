package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	switch os.Args[1] {
	case "hash":
		_, _ = hash(os.Args[2])
	case "compare":
		compare(os.Args[2], os.Args[3])
	}

}

func hash(password string) (string, error) {
	fmt.Println("Hash password...")
	output, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("error hashing: %v due to %s\n", password, err.Error())
		return "", err
	}

	hashOutput := string(output)
	fmt.Println("Returning hash: ", hashOutput)
	return hashOutput, nil
}

func compare(p, h string) bool {
	fmt.Print("Comparing ", p, " to ", h, " is ")
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	if err != nil {
		fmt.Println(err)
		return false
	}
	println("true")
	return true
}
