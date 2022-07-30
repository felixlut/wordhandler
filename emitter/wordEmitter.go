package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func fileToWordList(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var wordList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordList = append(wordList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return wordList, nil
}

type wordEmitter struct {
	wordList             []string
	host, port, connType string
	frequencyInMS, retryTime int
}

func (emitter wordEmitter) run() {
	target_host := emitter.host+":"+emitter.port
	fmt.Printf("Connecting to %s \n", target_host)
	retryAttempts := 0
	for retryAttempts < 10 {
		time.Sleep(time.Duration(emitter.frequencyInMS) * time.Millisecond)
		connection, err := net.Dial(emitter.connType, target_host)
		if err != nil {
			fmt.Printf("Failed to establish dial connection (%d attempts). Retry in %d seconds \n", retryAttempts, emitter.retryTime)
			fmt.Println(err)
			time.Sleep(time.Duration(emitter.retryTime) * time.Second)
			retryAttempts++
			
			continue
		}
		emitter.emitWord(connection)
		connection.Close()
		retryAttempts = 0
	}
}

func (emitter wordEmitter) emitWord(conn net.Conn) {
	word := emitter.wordList[rand.Intn(len(emitter.wordList))]
	_, err := conn.Write([]byte(word))
	if err != nil {
		fmt.Println("Failed to emit word")
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	wordList, err := fileToWordList("text.txt")
	if err != nil {
		fmt.Println("Failed to read file to list of words")
		panic(err)
	}

	emitter := wordEmitter{
		wordList:  wordList,
		host:      getEnv("TARGET_HOST", ""),
		port:      "9988",
		connType:  "tcp",
		frequencyInMS: 100,
		retryTime: 10,
	}

	emitter.run()
}
