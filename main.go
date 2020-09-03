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
	"sync"
	"time"
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

//producePassword is a function that generates random passwords
//and send them to password channel
func producePassword(password chan string, wg *sync.WaitGroup) {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012346789!#$%&'()*+,-./:;<=>?@[]^_`{|}~\"")

	b := make([]rune, rand.Intn(64-1)+1)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	password <- string(b)
	wg.Done()
}

//writePassword creates the file for if it does not exist
// and writes the password to it.
//It then reads the random numbers from the password channel and writes to the file
func writePassword(password chan string, done chan bool) {
	f, err := os.Create("output.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	for pass := range password {
		data, err := json.MarshalIndent(getPasswordReport([]byte(pass)), "", "\t")
		if err != nil {
			panic(err)
		}

		if _, err := f.Write(data); err != nil {
			println(err)
		}
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		done <- false
		return
	}
	done <- true
}

func main() {
	password := make(chan string)
	done := make(chan bool)

	//waitgroup
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go producePassword(password, &wg)
	}

	go writePassword(password, done)

	go func() {
		wg.Wait()
		close(password)
	}()

	d := <-done
	if d == true {
		fmt.Println("File written successfully")
	} else {
		fmt.Println("File writing failed")
	}
}
