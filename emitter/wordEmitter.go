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
	frequency, retryTime int
}

func (emitter wordEmitter) run() {
	retryAttempts := 0
	for retryAttempts < 10 {
		time.Sleep(time.Duration(emitter.frequency) * time.Second)
		connection, err := net.Dial(emitter.connType, emitter.host+":"+emitter.port)
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
	wordList, err := fileToWordList("text_short.txt")
	if err != nil {
		fmt.Println("Failed to read file to list of words")
		panic(err)
	}

	var host string
	switch env := os.Getenv("DEPLOY_ENVIRONMENT"); env {
	case "compose", "kubernetes":
		host = "receiver"
	default:
		host = ""
	}

	emitter := wordEmitter{
		wordList:  wordList,
		host:      host,
		port:      "9988",
		connType:  "tcp",
		frequency: 1,
		retryTime: 10,
	}

	emitter.run()
}
