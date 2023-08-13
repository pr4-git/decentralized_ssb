package gossip

import (
	"log"
	"net"
	"ssb-ng/dba"
	"ssb-ng/secretstream"

	"go.cryptoscope.co/netwrap"
)

func ConnectSyncClient(deps dba.Dependency, appKey []byte, serverPubKey []byte) error {
	client, err := secretstream.NewClient(*deps.GetKeyPair(), appKey)
	if err != nil {
		log.Printf("Cannnot join the gossip network. Details: %s", err)
		return err
	}

	if isAddrAvailable("127.0.0.1:8008") {
		log.Printf("Server Uninitiated. No gossip network to join.")
		return err
	}

	conn, err := netwrap.Dial(&net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8008,
	}, client.ConnWrapper(serverPubKey))

	if err != nil {
		log.Printf("Cannot establish connection to the gossip broker. Details: %s", err)
		return err
	}

	go handleConn(conn, deps)

	return err
}
