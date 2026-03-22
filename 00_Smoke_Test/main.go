package main;
import ("io"; "log"; "net");

func main(){
	listener, err := net.Listen("tcp", "0.0.0.0:8080");

	if err != nil {
		log.Fatal(err);
	}

	defer listener.Close();
	defer log.Println("Closing Server");

	log.Println("Server listening on port 8080");

	for {
		conn, err := listener.Accept();
		if err != nil {
			log.Println(err);
			continue;
		}

		log.Println("Client connection accepted, spawning a Go Routine");
		go handleClients(conn);
	}
}

func handleClients(conn net.Conn){
	defer conn.Close();
	defer log.Println("Closing Client Connection")

	_, err := io.Copy(conn, conn);
	if err != nil {
		log.Println(err);
	}
}

// imports

// start main function

	// create server using Listen command

	// handle error

	// defer the closage of server

	// start a loop and listen to clients

		// Accept connections

		// handle errors

		// create go routine and send the function that handles this client


// define function that handles each client
	// receive a connection object from main

	// read and write back

	// close when client sends eof