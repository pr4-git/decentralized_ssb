package securerpc

import (
	"itsy/ssb/dba"
	"itsy/ssb/handshake"
	"itsy/ssb/secretstream"
	"log"
	"net"
	"net/rpc"

	"github.com/jmoiron/sqlx"
	"go.cryptoscope.co/netwrap"
)

type Handler struct {
	DB *sqlx.DB
}

type Reply struct {
	Posts []dba.Post
}

func (rh *Handler) GetPosts(payload int, reply *Reply) error {
	posts, err := dba.FetchWall(rh.DB)
	*reply = Reply{Posts: posts}
	println("server done!")
	if err != nil {
		return err
	}
	return nil
}

func RunServer(db *sqlx.DB, keypair *handshake.EdKeyPair, appkey string) {
	server, err := secretstream.NewServer(*keypair, []byte(appkey))
	if err != nil {
		log.Printf("Server down for syncronization. %s", err)
	}
	listener, err := netwrap.Listen(
		&net.TCPAddr{
			IP:   net.IP{127, 0, 0, 1},
			Port: 8005,
		},
		server.ListenerWrapper(),
	)
	log.Println("You can connect with this server at:")
	log.Printf("%s@192.168.0.1:8005\n",
		server.Addr().String())
	if err != nil {
		log.Printf("Server down for syncronization. %s", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Cannot accept client. %s", err)
			}

			rpcServer := rpc.NewServer()
			// register rpc handers here
			h := Handler{DB: db}
			err = rpcServer.Register(&h)
			if err != nil {
				log.Printf("Syncronization down. %s", err)
			}
			rpcServer.ServeConn(conn)

			conn.Close()
		}
	}()
}