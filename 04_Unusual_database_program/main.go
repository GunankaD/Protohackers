package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

func main(){
	protocol, ip, port := "udp", "0.0.0.0", "8080"

	conn, err := net.ListenPacket(protocol, ip + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	db := &Database{
		db: make(map[string]string),
	}
	db.db["version"] = "Guna v1.6.7"

	buf := make([]byte, 1000)
	
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		req := string(buf[0:n])

		before, after, found := strings.Cut(req, "=")

		// lock mutex
		db.mu.Lock()

		if found { // insert
			if before != "version" {
				db.db[before] = after
			}
			db.mu.Unlock()
		} else { // retreive
			value, ok := db.db[before]
			db.mu.Unlock()

			if ok {
				res := fmt.Sprintf("%v=%v", before,value)
				conn.WriteTo([]byte(res), addr)
			}
		}
		
	}
}

type Database struct {
	mu sync.Mutex
	db map[string]string
}