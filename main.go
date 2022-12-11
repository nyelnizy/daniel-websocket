package main

import (
	dannyws "github.com/nyelnizy/dannyws/pkg"
)
type MyWsHandler struct{

}
func (h *MyWsHandler) OnOpen(c *dannyws.DannyWsConn)  {
	
}
func (h *MyWsHandler) OnMessage(m []byte,c *dannyws.DannyWsConn)  {
	
}
func (h *MyWsHandler) OnError(e error)  {
	
}
func (h *MyWsHandler) OnClose(e error)  {
	
}

func main() {
	h := &MyWsHandler{}
	ws := dannyws.WebsocketServer{
		Addr:   ":6001",
		WsHandler: h,
	}
	ws.Start()
}
