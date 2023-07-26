package main

import (
	"bufio"
	"fmt"
	"itsy/ssb/handshake"
	"itsy/ssb/secretstream"
	"log"
	"net"

	"go.cryptoscope.co/netwrap"
)

var appkey = []byte("SHSNet")

var (
	clientKeys, serverKeys *handshake.EdKeyPair
	appKey                 []byte
)

func init() {
	var err error
	clientKeys, err = handshake.GenEdKeyPair(nil)
	if err != nil {
		log.Fatalf("Unable to create keys for client. (Error: %v)", err)
	}
	serverKeys, err = handshake.GenEdKeyPair(nil)
	if err != nil {
		log.Fatalf("Unable to create keys for server. (Error: %v)", err)
	}

	appKP, err := handshake.GenEdKeyPair(nil)
	if err != nil {
		log.Fatalf("Unable to create key for application. (Error %v)", err)
	}

	appKey = appKP.Public
}

func main() {
	server, err := secretstream.NewServer(*serverKeys, appKey)
	if err != nil {
		log.Fatalf("Couldn't create server (Error: %v)", err)
	}

	listener, err := netwrap.Listen(
		&net.TCPAddr{
			IP:   net.IP{127, 0, 0, 1},
			Port: 8080},
		server.ListenerWrapper())
	if err != nil {
		log.Fatalf("netwrap couldn't create connection (Error: %v)", err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("couldn't accept conn (Error %v)", err)
		}
		defer conn.Close()

		for {
			connReader := bufio.NewReader(conn)
			msg, err := connReader.ReadBytes('\n')
			if err != nil {
				log.Fatalf("error reading from connection (Error: %v)", err)
			}
			fmt.Printf("client: %s", msg)
			/// -----------write:
			conn.Write([]byte("All your bases are belong to us!!\n"))
		}
	}()

	client, err := secretstream.NewClient(*clientKeys, appKey)
	if err != nil {
		log.Fatalf("couldn't create client (Error: %v)", err)
	}

	tcpAddr := netwrap.GetAddr(listener.Addr(), "tcp")
	connWrap := client.ConnWrapper(serverKeys.Public)

	conn, err := netwrap.Dial(tcpAddr, connWrap)
	if err != nil {
		log.Fatalf("couldn't dial server (Error: %v)", err)
	}
	defer conn.Close()

	for i := 0; i < 10; i++ {
		//--- write:
		conn.Write([]byte("oh naurrr!!!!!!\n"))
		// ------- read:
		connReader := bufio.NewReader(conn)
		msg, err := connReader.ReadBytes('\n')
		if err != nil {
			log.Fatalf("error reading from connection (Error: %v)", err)
		}
		fmt.Printf("server: %s", msg)
	}
}
