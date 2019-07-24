package nats

import (
	"github.com/jansemmelink/log"
	"github.com/jansemmelink/msvc"
)

//nats implements msvc.IServer to server micro-services from a NATS topic
type natsServer struct {
	URL string `json:"url" doc:"URL of NATS server. Defaults to \"localhost:4222\""`
}

func (ns *natsServer) Validate() error {
	if len(ns.URL) == 0 {
		ns.URL = "localhost:4222"
	}
	log.Debugf("Validated %T", ns)
	return nil
}

func init() {
	msvc.RegisterServer("nats", &natsServer{})
}
