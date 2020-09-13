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
	"math/rand"
	"os"
	"time"
)

var (
	charactersList = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012346789!#$%&'()*+,-./:;<=>?@[]^_`{|}~\""
	charactersCount = 64
	outputFile = "output.json"
)

type PasswordReport struct {
	MD5    string `json:"md5"`
	SHA1   string `json:"sha-1"`
	SHA256 string `json:"sha-256"`
	SHA512 string `json:"sha-512"`
}

type PasswordMap struct{
	Key 	 PasswordReport
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

func checkForDuplicate(filepath string, password string) bool{

	jsonFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	bytesData, _ := ioutil.ReadAll(jsonFile)
	firstJSON := []byte("[")
	lastJSONBracket := []byte("]")
	bytesData = append(firstJSON, bytesData...)
	bytesData = append(bytesData, lastJSONBracket...)

	var passwords []map[string]PasswordReport
	err = json.Unmarshal(bytesData, &passwords)
	if err != nil{
		panic(err)
	}

	//get all json keys
	keys := make(map[string]bool)
	for _, singleMap := range passwords{
		for k := range singleMap {
			keys[k]=true
		}
	}

	if val, ok := keys[password]; ok {
		return val
	}

	return false
}

func main() {
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("OS Error: %s", err)
	}

	for {
		randomPass := RandomString()
		b, err := json.Marshal(getCompleteReport([]string{(randomPass)}))
		if err != nil {
			panic(err)
		}
		exists:= checkForDuplicate(outputFile, randomPass)

		if !exists {
			fileStats, err := file.Stat()
			if err != nil {
				panic(err)
			}
			if fileStats.Size() != 0 {
				b = append([]byte(","), b...)
			}
			if _, err = file.Write(b); err != nil {
				fmt.Printf("Error writing to a file %s", err)
			}
		}
	}
}
