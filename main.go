package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"encoding/json"
)

/* Todo:
- Generate a random string (Duplicate Check)
- Get the {SHA-512|SHA-256|SHA-1|MD5} of the string
- Save the {SHA-512|SHA-256|SHA-1|MD5} with the string in a json file.
*/

func GenerateString() {
	// Generate a random string
}

func GetHash() {
	aStringToHash := []byte("Lets hash this sentence!")

	//Get the hashes in bytes
	md5Bytes := md5.Sum(aStringToHash)
	sha1Bytes := sha1.Sum(aStringToHash)
	sha256Bytes := sha256.Sum256(aStringToHash)
	sha512Bytes := sha512.Sum512(aStringToHash)

	//Print out what will be hashed
	fmt.Println(string(aStringToHash))

	//Bytes to string
	fmt.Println("MD5:", hex.EncodeToString(md5Bytes[:]))
	fmt.Println("SHA-1:", hex.EncodeToString(sha1Bytes[:]))
	fmt.Println("SHA-256:", hex.EncodeToString(sha256Bytes[:]))
	fmt.Println("SHA-512:", hex.EncodeToString(sha512Bytes[:]))
}

func main() {
	// save the output to a file.
}
