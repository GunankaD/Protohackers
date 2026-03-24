package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

func main(){
	protocol, ip, port := "tcp", "0.0.0.0", "8080"
	listener, err := net.Listen(protocol, ip + ":" + port)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	log.Printf("Server listening on port: %v\n", port)

	chatRoom := &Room{
		users: make(map[net.Conn]string),
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleClient(conn, chatRoom)
	}
}

func handleClient(conn net.Conn, r *Room) {
	defer conn.Close()

	log.Println("Handling a new Client")

	// 1. Ask their name
	welcomeMsg := "Welcome to budgetchat! What shall I call you?\n"
	conn.Write([]byte(welcomeMsg))

	// 2. Read their name
	scanner := bufio.NewScanner(conn)
	var name string

	if scanner.Scan() {
		name = scanner.Text()
		if !isValidName(name) {
			return
		}

		r.AddUser(name, conn)
	} else { // if they leave before setting their name, silently cut the connection
		return
	}

	// 3. Handle their broadcast messages
	for scanner.Scan() {
		msg := scanner.Text()
		broadcastMsg := fmt.Sprintf("[%v] %v\n", name, msg);
		r.Broadcast(conn, broadcastMsg)
	}

	// 4. Remove user when they leave
	r.RemoveUser(conn)
}

func isValidName(name string) bool {
	if len(name) == 0 {
		return false
	}

	name = strings.ToLower(name)
	for _, char := range name {
		if !(('a' <= char && char <= 'z') || ('0' <= char && char <= '9')) {
			return false
		}
	}
	return true
}

type Room struct {
	mu sync.Mutex
	users map[net.Conn]string
}

func (r *Room) AddUser(name string, conn net.Conn){
	// 1. Lock the mutex
	r.mu.Lock()

	// 2. Fetch all current users
	var currentUsers []string
	for _, userName := range r.users {
		currentUsers = append(currentUsers, userName)
	}

	// 3. Add the new user
	r.users[conn] = name

	// 4. Prepare current users message for the new user and Write back
	currentUsersMsg := "* The room contains:"
	for index, userName := range currentUsers {
		if index == 0 {
			currentUsersMsg += " " + userName
		} else {
			currentUsersMsg += ", " + userName
		}
	}
	currentUsersMsg += "\n"
	conn.Write([]byte(currentUsersMsg))

	// 5. Release the lock
	r.mu.Unlock()

	// 6. Broadcast
	msg := fmt.Sprintf("* a wild %v appeared!\n", name)
	r.Broadcast(conn, msg)
}

func (r *Room) RemoveUser(conn net.Conn){
	// 1. Lock the mutex
	r.mu.Lock()

	// 2. Note user's name and remove the user
	name := r.users[conn]
	delete(r.users, conn)

	// 3. Unlock the mutex
	r.mu.Unlock()

	// 4. Prepare broadcast message and call the function
	msg := fmt.Sprintf("* %v has left the chat!\n", name)
	r.Broadcast(conn, msg)
}

func (r *Room) Broadcast(senderConn net.Conn, msg string){
	// 1. Lock the mutex
	r.mu.Lock()

	// 2. Collect all current users
	var currentUsers []net.Conn
	for conn := range r.users {
		if conn != senderConn {
			currentUsers = append(currentUsers, conn)
		}
	}

	// 3. Release Mutex
	r.mu.Unlock()

	// 4. Broadcast to everyone present in the network
	for _, conn := range currentUsers {
		conn.Write([]byte(msg))
	}
}