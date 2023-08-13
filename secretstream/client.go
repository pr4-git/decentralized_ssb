package secretstream

import (
	"fmt"
	"net"
	"ssb-ng/boxstream"
	"ssb-ng/handshake"
	"time"

	"go.cryptoscope.co/netwrap"
)

// Client can dial secret-handshake server endpoints
type Client struct {
	appKey []byte
	kp     handshake.EdKeyPair
}

// NewClient creates a new Client with the passed keyPair and appKey
func NewClient(kp handshake.EdKeyPair, appKey []byte) (*Client, error) {
	// TODO: consistancy check?!..
	return &Client{
		appKey: appKey,
		kp:     kp,
	}, nil
}

// ConnWrapper returns a connection wrapper for the client.
func (c *Client) ConnWrapper(pubKey []byte) netwrap.ConnWrapper {
	return func(conn net.Conn) (net.Conn, error) {
		state, err := handshake.NewClientState(c.appKey, c.kp, pubKey)
		if err != nil {
			return nil, err
		}

		errc := make(chan error)
		go func() {
			errc <- handshake.Client(state, conn)
			close(errc)
		}()

		select {
		case err := <-errc:
			if err != nil {
				return nil, err
			}
		case <-time.After(30 * time.Second):
			return nil, fmt.Errorf("secretstream: handshake timeout")
		}

		enKey, enNonce := state.GetBoxstreamEncKeys()
		deKey, deNonce := state.GetBoxstreamDecKeys()

		boxed := &Conn{
			boxer:   boxstream.NewBoxer(conn, &enNonce, &enKey),
			unboxer: boxstream.NewUnboxer(conn, &deNonce, &deKey),
			conn:    conn,
			local:   c.kp.Public[:],
			remote:  state.Remote(),
		}

		return boxed, nil
	}
}
