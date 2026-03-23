package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"log"
)

func main() {
	// 1. Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error connecting:", err)
	}
	defer conn.Close()

	// 2. Define our test data
	inserts := [][2]int32{
		{12345, 101},
		{12346, 102},
		{12347, 100},
		{40960, 5},
	}

	queries := [][2]int32{
		{12988, 40960},
	}

	// 3. Send Inserts
	for _, ins := range inserts {
		// Prepare a 9-byte message: 'I' (1 byte) + Timestamp (4) + Price (4)
		binary.Write(conn, binary.BigEndian, byte('I'))
		binary.Write(conn, binary.BigEndian, ins[0])
		binary.Write(conn, binary.BigEndian, ins[1])
	}

	// 4. Send Queries
	for _, qs := range queries {
		binary.Write(conn, binary.BigEndian, byte('Q'))
		binary.Write(conn, binary.BigEndian, qs[0])
		binary.Write(conn, binary.BigEndian, qs[1])
	}
	
    
	// 5. Read the 4-byte response
	var result int32
	err = binary.Read(conn, binary.BigEndian, &result)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Printf("Received Mean: %d\n", result)
}