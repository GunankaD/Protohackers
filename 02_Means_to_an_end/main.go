package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

func main(){
	protocol, ip, port := "tcp", "0.0.0.0", "8080"
	listener, err := net.Listen(protocol, ip + ":" + port)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	log.Println("Server listening on ", ip, ":", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleClients(conn)
	}
}


type History struct {
	Timestamp int32
	Price int32
}

func handleClients(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 9)
	var history []History

	for {
		_, err := io.ReadFull(conn, buffer)
		if err != nil {
			return
		}

		var queryType string
		var first,  second int32

		queryType = string(buffer[0])
		lerr := binary.Read(bytes.NewReader(buffer[1:5]), binary.BigEndian, &first)
		rerr := binary.Read(bytes.NewReader(buffer[5:9]), binary.BigEndian, &second)

		if lerr != nil || rerr != nil {
			return
		}

		if queryType == "I" {
			newData := History{
				Timestamp : first,
				Price : second,
			}
			history = append(history, newData)
		} else {
			var mean, count, sum int64

			for _, data := range history {
				if first <= data.Timestamp && data.Timestamp <= second {
					sum += int64(data.Price)
					count ++
				}
			}

			if count != 0 {
				mean = sum / count
				binary.Write(conn, binary.BigEndian, int32(mean))
			} else {
				binary.Write(conn, binary.BigEndian, int32(count))
			}
		}
		
	}
}