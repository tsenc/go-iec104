package asdu

import (
	"net"
)

// Connect interface
type Connect interface {
	ServerId() string
	Params() *Params
	Send(a *ASDU) error
	UnderlyingConn() net.Conn
}
