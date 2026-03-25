package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
	"regexp"
)

var boguscoinRegex = regexp.MustCompile(`^7[a-zA-Z0-9]{25,34}$`)

func main(){
	myProtocol, myIp, myPort := "tcp", "0.0.0.0", "8080"

	downListener, err := net.Listen(myProtocol, myIp + ":" + myPort)
	if err != nil {
		log.Fatal(err)
	}
	defer downListener.Close()

	log.Println("Middleman listening on port: ", myPort)

	for {
		downConn, err := downListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleClient(downConn)
	}

}

func handleClient(downConn net.Conn){
	defer downConn.Close()

	upProtocol, upIp, upPort := "tcp", "chat.protohackers.com", "16962"
	upConn, err := net.Dial(upProtocol, upIp + ":" + upPort)
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go repeater(upConn, downConn, &wg) // client -> us -> server
	go repeater(downConn, upConn, &wg) // server -> us -> client

	wg.Wait()
}

func repeater(dest net.Conn, src net.Conn, wg *sync.WaitGroup){
	defer wg.Done()
	defer dest.Close()
	defer src.Close()

	scanner := bufio.NewScanner(src)
	tonysAddress := "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

	for scanner.Scan() {
		req := scanner.Text()

		// search and replace the boguscoin
		words := strings.Split(req, " ")

		for index, word := range words {
			if boguscoinRegex.MatchString(word) {
				words[index] = tonysAddress
			}
		}

		modifiedRes := strings.Join(words, " ")
		dest.Write([]byte(modifiedRes + "\n"))
	}
}