package gossip

import (
	"encoding/base64"
	"encoding/gob"
	"io"
	"log"
	"net"
	"ssb-ng/dba"
	"ssb-ng/secretstream"
	"strings"
	"time"

	"go.cryptoscope.co/netwrap"
)

var IsServer bool = false

func StartSyncServer(deps dba.Dependency, appkey []byte) {
	server, err := secretstream.NewServer(*deps.GetKeyPair(), appkey)
	if err != nil {
		log.Printf("Cannot initialize the gossip server. Err %s", err)
		return
	}
	if !isAddrAvailable("127.0.0.1:8008") {
		log.Printf("Cannot start instance as gossip server.")
		return
	}

	listener, err := netwrap.Listen(&net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8008,
	}, server.ListenerWrapper())
	if err != nil {
		log.Println("Cannot listen on addr 127.0.0.1:80085. Quitting gossip server.")
		return
	}
	defer listener.Close()

	log.Printf("🔥🔥Sync server listening at: %s", listener.Addr())
	IsServer = true

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn, deps)
	}

}

func handleConn(conn net.Conn, deps dba.Dependency) {

	defer conn.Close()

	postencoder := gob.NewEncoder(conn)
	postdecoder := gob.NewDecoder(conn)

	for range time.Tick(1 * time.Second) {
		// changed, err := dba.RecordsChanged(deps)
		// if err != nil {
		// 	log.Printf("handleConn err: %s", err)
		// 	continue
		// }
		changed := true
		if changed {
			posts, err := dba.FetchAllPosts(deps)
			if err != nil {
				log.Printf("handleConn err: %s", err)
				continue
			}

			// recv posts
			go func() {
				for {
					var post dba.Post
					err := postdecoder.Decode(&post)
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Printf("Decoder errored. Details: %s", err)
					}
					dba.InsertPostToDB(deps, &post)
				}
			}()

			// send posts
			for i := range posts {
				err := postencoder.Encode(posts[i])
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("Encoder erroered. Details: %s", err)
					continue
				}
			}

		}
	}

}

func isAddrAvailable(addr string) bool {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func addrToPk(addr string) ([]byte, error) {
	pkIdentifier := strings.Split(addr, "|@")[1]
	pkeyb64 := strings.Split(pkIdentifier, ".")[0]
	return base64.StdEncoding.DecodeString(pkeyb64)
}
