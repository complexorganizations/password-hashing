package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

/* Todo:
- Make it write on new lines
- Add Concurrency
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
	var found bool
	for _, pass := range passwords {
		if _, found = completeReport[pass]; !found {
			completeReport[pass] = getPasswordReport([]byte(pass))
		}
	}
	return completeReport
}

func RandomString() string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012346789!#$%&'()*+,-./:;<=>?@[]^_`{|}~\"")

	b := make([]rune, rand.Intn(64-1)+1)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	file, err := os.Create("output.json")
	if err != nil {
		fmt.Printf("os erreur %s", err)
	}

	for {
		// replace with Marshal, im using MarshalIndent for nice formatting
		b, err := json.MarshalIndent(getCompleteReport([]string{(RandomString())}), " ", "\t")
		if err != nil {
			panic(err)
		}
		//ioutil.WriteFile("output.json", data, 0644)
		if _, err = file.Write(b); err != nil {
			fmt.Printf("error writing to a file %s", err)
		}
	}
}
