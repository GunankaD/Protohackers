package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math"
	"net"
)

func main(){
	protocol, ip, port := "tcp", "0.0.0.0", "8080"
	listener, err := net.Listen(protocol, ip + ":" + port)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	log.Println("Server now listening on " + ip + ":" + port)

	// Create new Go Routine for each client
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleClients(conn)
	}
}

type Request struct {
	Method *string `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime bool `json:"prime"`
}

func handleClients(conn net.Conn){
	defer conn.Close()
	defer log.Println("Closing a client connection")
	log.Println("Serving a new client")

	scanner := bufio.NewScanner(conn)
	noOfRequests := 0

	for scanner.Scan() {
		noOfRequests ++
		log.Printf("Processing request number: %v\n", noOfRequests)

		var requestStruct Request
		var responseStruct Response
		var responseBytes []byte

		requestString := scanner.Text()
		log.Printf("Client's request: %v\n", requestString)
	
		err := json.Unmarshal([]byte(requestString), &requestStruct)

		// error with unmarshalling or missing of required fields
		if err != nil || requestStruct.Method == nil ||  requestStruct.Number == nil || *requestStruct.Method != "isPrime"  {
			log.Println("Malformed request, exiting...")

			responseStruct.Method = "malformedResponse"
			responseStruct.Prime = false

			responseBytes, _ = json.Marshal(responseStruct)

			conn.Write(responseBytes)
			conn.Write([]byte("\n"))
			return // end the connection
		}

		responseStruct.Method = "isPrime"
		responseStruct.Prime = isPrime(*requestStruct.Number)

		responseBytes, _ = json.Marshal(responseStruct)
		log.Printf("Conforming request. Server response: %v\n", string(responseBytes))

		conn.Write(responseBytes)
		conn.Write([]byte("\n"))
	}

}

func isPrime(num float64) bool {

	// primality is not defined over floats
	if math.Floor(num) != math.Ceil(num) {
		return false
	}

	x := int(num)

	// handle negative numbers, 0s and 1s
	if x <= 1 { 
		return false
	}

	for i := 2; i * i <= x; i ++ {
		if x % i == 0 {
			return false
		}
	}

	return true
}