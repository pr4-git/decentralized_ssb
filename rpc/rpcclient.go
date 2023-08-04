package securerpc

import (
	"encoding/base64"
	"itsy/ssb/handshake"
	"itsy/ssb/secretstream"
	"log"
	"net"
	"net/rpc"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.cryptoscope.co/netwrap"
)

func SyncClient(db *sqlx.DB, keypair *handshake.EdKeyPair, serveraddr string, appkey string) {
	splitstr := strings.Split(serveraddr, "|@")
	log.Printf("%v", splitstr)
	netaddr, serverkeystr := splitstr[0], strings.TrimSuffix(splitstr[1], ".ed25519")
	log.Printf("%s,%s", netaddr, serveraddr)
	serverKey, err := base64.StdEncoding.DecodeString(serverkeystr)
	if err != nil {
		log.Printf("Cannot syncronize with the server %s", err)
		return
	}

	client, err := secretstream.NewClient(*keypair, []byte(appkey))
	if err != nil {
		log.Fatalf("couldn't create client (Error: %v)", err)
	}

	serverAddr, err := net.ResolveTCPAddr("tcp", netaddr)
	if err != nil {
		log.Printf("Cannot syncronize with the server %s", err)
		return
	}
	tcpAddr := netwrap.GetAddr(
		serverAddr,
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

	for i := 0; i < len(reply.Posts); i++ {
		err := reply.Posts[i].SyncToFeed(db)
		if err != nil {
			log.Printf("Possible spoof attempt. %s", err)
			break
		}
	}
	log.Printf("Sync complete with %s", base64.StdEncoding.EncodeToString(serverKey))
}
