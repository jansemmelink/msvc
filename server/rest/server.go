package rest

import (
	"github.com/jansemmelink/log"
	"github.com/jansemmelink/msvc"
)

//restServer implements msvc.IServer to be a HTTP REST interface for micro-services
type restServer struct {
}

func (rs restServer) Validate() error {
	log.Debugf("Validated %T", rs)
	return nil
}

func init() {
	msvc.RegisterServer("rest", restServer{})
}
