package main

import (
	"fmt"
	"itsy/ssb/handshake"
	"log"
	"net"
	"os"
)

var appkey = []byte("SHSNet")

func main() {
	serverKey, err := handshake.GenEdKeyPair(nil)
	if err != nil {
		log.Fatal(err)
	}

	clientKey, err := handshake.GenEdKeyPair(nil)
	if err != nil {
		log.Fatal(err)
	}

	serverState, err := handshake.NewServerState(appkey, *serverKey)
	if err != nil {
		log.Fatal(err)
	}
	clientState, err := handshake.NewClientState(appkey, *clientKey, serverKey.Public)
	if err != nil {
		log.Fatal(err)
	}

	server, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	clientconn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	serverconn, err := server.Accept()
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("client connected")
	// server
	go func() {
		err := handshake.Server(serverState, serverconn)
		if err != nil {
			fmt.Println(err)
		}
		serverconn.Close()
	}()

	// client
	err = handshake.Client(clientState, clientconn)
	if err != nil {
		fmt.Println(err)
	}
	clientconn.Close()
}
