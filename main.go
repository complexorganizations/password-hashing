package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

/* Todo:
- Generate a random string (Duplicate Check)
- Get the {SHA-512|SHA-256|SHA-1|MD5} of the string
- Save the {SHA-512|SHA-256|SHA-1|MD5} with the string in a json file.
*/

type PasswordReport struct {
	MD5    string `json:"md5"`
	SHA1   string `json:"sha-1"`
	SHA256 string `json:"sha-256"`
	SHA512 string `json:"sha-512"`
}

func getPasswordReport(password []byte) *PasswordReport {
	md5Bytes := md5.Sum(password)
	sha1Bytes := sha1.Sum(password)
	sha256Bytes := sha256.Sum256(password)
	sha512Bytes := sha512.Sum512(password)

	return &PasswordReport{
		MD5:    hex.EncodeToString(md5Bytes[:]),
		SHA1:   hex.EncodeToString(sha1Bytes[:]),
		SHA256: hex.EncodeToString(sha256Bytes[:]),
		SHA512: hex.EncodeToString(sha512Bytes[:]),
	}
}

func getCompleteReport(passwords []string) map[string]*PasswordReport {
	completeReport := make(map[string]*PasswordReport)

	for _, pass := range passwords {
		completeReport[pass] = getPasswordReport([]byte(pass))
	}

	return completeReport
}

func main() {
	// replace with Marshal, im using MarshalIndent for nice formatting
	data, err := json.MarshalIndent(getCompleteReport([]string{"password", "abc123"}), "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
