package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func setupListener(connectionType, port string) (net.Listener, error) {
	connectionURL := ":" + port
	listener, err := net.Listen(connectionType, connectionURL)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		fmt.Printf("url: %s", connectionURL)
		return nil, err
	}
	fmt.Println("Listening on " + connectionURL)
	fmt.Println("Waiting for client...")

	return listener, nil
}

func readFromConnection(connection net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		return "", err
	}
	message := string(buffer[:mLen])
	return message, nil
}

type wordStat struct {
	word                  string
	firstSeen, lastSeen   time.Time
	timesSeen, sinceFlush int
}

func (stat wordStat) String() string {
	return fmt.Sprintf("Word: %s \nFirst seen: %s \nLast Seen: %s \nTimes seen: %d \nSince flush: %d \n",
		stat.word, stat.firstSeen.Format(time.Kitchen), stat.lastSeen.Format(time.Kitchen), stat.timesSeen, stat.sinceFlush)
}

type wordReceiver struct {
	wordStats                 map[string]wordStat
	mu                        sync.Mutex
	words                     []string
	port, cliPort, serverType string
	flushFrequency, retryTime int
}

func (receiver *wordReceiver) Value(key string) (wordStat, bool) {
	receiver.mu.Lock()
	defer receiver.mu.Unlock()
	value, ok := receiver.wordStats[key]
	return value, ok
}

func (receiver *wordReceiver) PutValue(key string, newValue wordStat) {
	receiver.mu.Lock()
	receiver.wordStats[key] = newValue
	receiver.mu.Unlock()
}

func (receiver *wordReceiver) handleCliCommand(connection net.Conn) {
	word, err := readFromConnection(connection)
	check(err)

	var response string
	if stat, ok := receiver.Value(word); ok {
		response = stat.String()
	} else {
		response = fmt.Sprintf("Word %s has not been seen", word)
	}

	if _, err := connection.Write([]byte(response)); err != nil {
		fmt.Println("Failed to send response to CLI")
		check(err)
	}
}

func (receiver *wordReceiver) runCliServer(wg *sync.WaitGroup) {
	defer wg.Done()

	server, err := setupListener(receiver.serverType, receiver.cliPort)
	check(err)
	defer server.Close()

	for {
		connection, err := server.Accept()
		check(err)
		go receiver.handleCliCommand(connection)
	}

}

func (receiver *wordReceiver) runWordServer(wg *sync.WaitGroup) {
	defer wg.Done()
	retryAttempts := 0
	var server net.Listener
	for retryAttempts < 10 {
		var err error
		server, err = setupListener(receiver.serverType, receiver.port)
		if err != nil {
			retryAttempts++
			fmt.Printf("Failed to establish listener connection (%d attempts). Retry in %d seconds \n", retryAttempts, receiver.retryTime)
			fmt.Println(err)
			time.Sleep(time.Duration(receiver.retryTime) * time.Second)
			continue
		}

		break
	}
	defer server.Close()

	// Continously catch and handle the sent words
	for {
		connection, err := server.Accept()
		check(err)
		go receiver.catchWord(connection)
	}
}

func (receiver *wordReceiver) catchWord(connection net.Conn) {
	word, err := readFromConnection(connection)
	check(err)

	if val, ok := receiver.Value(word); ok {
		val.lastSeen = time.Now()
		val.timesSeen++
		val.sinceFlush++
		receiver.PutValue(word, val)
	} else {
		val = wordStat{
			word:       word,
			lastSeen:   time.Now(),
			firstSeen:  time.Now(),
			timesSeen:  1,
			sinceFlush: 1,
		}
		receiver.PutValue(word, val)
		receiver.words = append(receiver.words, word)
	}
}

func (receiver *wordReceiver) runFlusher(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		time.Sleep(time.Duration(receiver.flushFrequency) * time.Second)
		for key, val := range receiver.wordStats {
			val.sinceFlush = 0
			receiver.wordStats[key] = val
		}
	}
}

func (receiver *wordReceiver) run() {
	var wg sync.WaitGroup
	wg.Add(3)

	go receiver.runWordServer(&wg)
	go receiver.runCliServer(&wg)

	go receiver.runFlusher(&wg)

	wg.Wait()
}

func main() {
	receiver := wordReceiver{
		wordStats:      make(map[string]wordStat),
		port:           "9988",
		cliPort:        "8899",
		serverType:     "tcp",
		flushFrequency: 10,
		retryTime:      10,
	}

	receiver.run()

}
