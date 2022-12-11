package dannyws

type WsHandler interface{
	OnOpen(*DannyWsConn)
	OnMessage([]byte,*DannyWsConn)
	OnClose(error)
	OnError(error)
}