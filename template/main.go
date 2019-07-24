//Package main is a template micro-service
package main

import (
	//other libraries
	"github.com/jansemmelink/log"

	//config sources that may be used:
	_ "github.com/jansemmelink/config/source/files"

	//micro-server server implementations that may be used:
	_ "github.com/jansemmelink/msvc/server/nats"
	_ "github.com/jansemmelink/msvc/server/rest"
)

func main() {
	//just to demonstrate how it works...
	log.DebugOn()

	//create the micro-service definition
	t := Template()

	//not necessary - just to demonstrate
	t.Test("hello", "{\"name\":\"Jan\"}")

	//serve on all configured interfaces
	t.Serve()
}
