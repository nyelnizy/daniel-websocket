package dannyws

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type connContext string

const (
	uuid                   = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	contextKey connContext = "con"
)

func saveInContext(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, contextKey, c)
}

type WebsocketServer struct {
	Addr             string
	Origins          []string
	SubProtocols     []string
	SelectedProtocol string
	WsHandler        WsHandler
}
type DannyWsConn struct {
	conn *net.Conn
}

func (d *DannyWsConn) Send(message []byte) {

}
func (d *DannyWsConn) Close() {

}

func (s *WebsocketServer) Start() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("New request...")
		e := ""
		if r.Method != http.MethodGet {
			e = "method not supported"
			handleError(s, e)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(e))
			return
		}
		// if origins is empty then all origins are supported
		exists := len(s.Origins) == 0
		for _, o := range s.Origins {
			if o == r.Header.Get("Origin") {
				exists = true
				break
			}
		}
		if !exists {
			e = "origin not supported"
			handleError(s, e)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(e))
			return
		}
		// compute and encode key for client to verify
		key := strings.Trim(r.Header.Get("Sec-WebSocket-Key"), " ")

		if key == "" {
			e = "client is not supported"
			handleError(s, e)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(e))
			return
		}
		concatenatedString := strings.Builder{}
		concatenatedString.Write([]byte(key))
		concatenatedString.Write([]byte(uuid))

		h := sha1.New()
		h.Write([]byte(concatenatedString.String()))
		hash := h.Sum(nil)
		secKey := base64.StdEncoding.EncodeToString(hash)

		conn := r.Context().Value(contextKey).(net.Conn)

		status := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n", secKey)
		conn.Write([]byte(status))
		// spin up a goroutine for current connection
		go s.handleConnection(&conn)
	})
	server := &http.Server{
		Addr:        s.Addr,
		ConnContext: saveInContext,
	}
	fmt.Println("Starting Server....")
	server.ListenAndServe()
}

func handleError(s *WebsocketServer, errorMessage string) {
	err := errors.New(errorMessage)
	if s.WsHandler != nil {
		s.WsHandler.OnError(err)
		s.WsHandler.OnClose(err)
	}
}

func (s *WebsocketServer) handleConnection(conn *net.Conn) {
	c := &DannyWsConn{
		conn: conn,
	}
	if s.WsHandler != nil {
		s.WsHandler.OnOpen(c)
	}
}
