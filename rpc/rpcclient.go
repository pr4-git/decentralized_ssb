package securerpc

import (
	"encoding/base64"
	"itsy/ssb/dba"
	"itsy/ssb/handshake"
	"itsy/ssb/secretstream"
	"log"
	"net"
	"net/rpc"
	"sort"

	"github.com/jmoiron/sqlx"
	"go.cryptoscope.co/netwrap"
	"golang.org/x/crypto/ed25519"
)

func SyncClient(db *sqlx.DB, keypair *handshake.EdKeyPair, serverKey ed25519.PublicKey, appkey string) {
	client, err := secretstream.NewClient(*keypair, []byte(appkey))
	if err != nil {
		log.Fatalf("couldn't create client (Error: %v)", err)
	}

	tcpAddr := netwrap.GetAddr(
		&net.TCPAddr{
			IP:   net.IP{127, 0, 0, 1},
			Port: 8005,
		},
		"tcp")
	connWrap := client.ConnWrapper(serverKey)

	conn, err := netwrap.Dial(tcpAddr, connWrap)
	if err != nil {
		log.Fatalf("couldn't dial server (Error: %v)", err)
	}
	defer conn.Close()
	log.Println("Connected!!!")

	clientRPC := rpc.NewClient(conn)

	var reply Reply
	err = clientRPC.Call("Handler.GetPosts", 0, &reply)
	if err != nil {
		log.Printf("Error in syncronization. %s", err)
	}

	sort.Slice(reply.Posts, func(i, j int) bool {
		return reply.Posts[i].ID < reply.Posts[j].ID
	})

	err = dba.NewProfile(serverKey, "").FollowProfile(db)
	if err != nil {
		log.Printf("Unable to create a peer. %s", err)
	}

	for i := 0; i < len(reply.Posts); i++ {
		err := reply.Posts[i].SyncToFeed(db)
		if err != nil {
			log.Printf("Possible spoof attempt. %s", err)
			break
		}
	}
	log.Printf("Sync complete with %s", base64.StdEncoding.EncodeToString(serverKey))
}
