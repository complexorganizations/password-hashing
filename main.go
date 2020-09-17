package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"
)

var (
	charactersList  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012346789!"#$%&'()*+,-./:;<=>?@[\]^_`{|}~"
	charactersCount = 64
	outputFile      = "output.json"
)

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
	var letterRunes = []rune(charactersList)

	b := make([]rune, rand.Intn(charactersCount)+1)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func appendIfNotPresent(pr map[string]*PasswordReport) {
	result := []map[string]*PasswordReport{}
	//check if file exits
	if !fileExists(outputFile) {
		//simple append
		result = append(result, pr)
		b, err := json.Marshal(&result)
		if err != nil {
			log.Fatalln("Error while marshaling: ", err)
		}
		writingInFile(b)
		return
	}

	//check for duplications
	byteValue, err := ioutil.ReadFile(outputFile)
	if err != nil {
		log.Fatalln("Error while reading: ", err)
	}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		log.Fatalln("Error while marshaling: ", err)
	}

	isFound := false
	for _, r := range result {
		if reflect.DeepEqual(r, pr) {
			isFound = true
			break
		}
	}

	if isFound {
		fmt.Println("Duplicate Found.")
		return
	}

	//appending
	result = append(result, pr)
	b, err := json.Marshal(&result)
	if err != nil {
		log.Fatalln("Error while marshaling: ", err)
	}
	writingInFile(b)
}

func writingInFile(b []byte) {
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("OS Error: %s", err)
	}
	if _, err = file.Write(b); err != nil {
		fmt.Printf("Error writing to a file %s", err)
	}
	file.Close()
}

func main() {
	for {
		appendIfNotPresent(getCompleteReport([]string{RandomString()}))
	}
}
