package msvc

import (
	"sync"

	"github.com/jansemmelink/config"
	"github.com/jansemmelink/log"
)

//IServer is a micro-service interface to process requests
//it could be HTTP REST with application/json etc...
type IServer interface {
	//server must be configurable, so embed this:
	config.IValidator
}

//RegisterServer must be called in the server implementation's init() func
//to make it available to the micro-service framework. It will be constructed
//if the server name is configured
func RegisterServer(name string, tmpl IServer) {
	if len(name) == 0 || tmpl == nil {
		panic("Server registration must have a name")
	}

	serverMutex.Lock()
	defer serverMutex.Unlock()

	if _, ok := serverTmpl[name]; ok {
		panic("Duplicate server name")
	}
	serverTmpl[name] = tmpl
}

var (
	serverMutex sync.Mutex
	serverTmpl  = make(map[string]IServer)
)

func startServers(cs config.ISet) {
	log.Debugf("Trying %d server configurations...", len(serverTmpl))
	for serverName, tmpl := range serverTmpl {
		log.Debugf("Trying %s...", serverName)
		serverConfiguration, err := cs.Add(serverName, tmpl)
		if err != nil {
			log.Debugf("server[%s] not configured: %+v", serverName, err)
			continue
		}

		configuredServer := serverConfiguration.Current().(IServer)
		log.Debugf("Got %s: %T: %+v", serverName, configuredServer, configuredServer)
	}

	panic("NYI")
}
