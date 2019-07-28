//Package nats implements a msvc.IServer to serve an IMicroService on a NATS subscription.
//That means the process will subscribe to NATS topic "<name>.*"
//and you send requests to that topic to have them served,
//for example by using the github.com:nats.io/examples/nats-req utility like this:
//
//  $ nats-req template.hello '{"request":{"value":123}}'
//	Published [template.hello] : '{"request":{"value":123}}'
//	Received  [_INBOX.C07XQjpNIGccsfAU3QDW6c.aSkIsJFD] : '{"header":{...}, "request":{...}, "result":{...}, "response":{...}}'
//
//One can also submit a header with constraints in the request, often with timeout value.
//When you receive will contain optional items for header, request, result and response.
package nats

import (
	"time"
	"strings"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/jansemmelink/log"
	"github.com/jansemmelink/msvc"
)

//nats implements msvc.IServer to server micro-services from a NATS topic
type natsServer struct {
	URL string `json:"url" doc:"URL of NATS server. Defaults to \"localhost:4222\""`

	//run-time private data:
	msvc msvc.IMicroService
}

func (ns *natsServer) Validate() error {
	if len(ns.URL) == 0 {
		ns.URL = "localhost:4222"
	}
	log.Debugf("Validated %T", ns)
	return nil
}

func (ns natsServer) Run(msvc msvc.IMicroService) {
	ns.msvc = msvc

	//connect to NATS
	conn, err := nats.Connect(ns.URL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second*2),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			log.Debugf("Trying to reconnect %+v\n", conn)
		}))
	if err != nil {
		panic(log.Wrapf(err, "Failed to connect to NATS server %s", ns.URL))
	}

	//make a queue subscription to start consuming messages from the topic
	/*subscription*/_, err = conn.QueueSubscribe(
		msvc.Name()+".*",
		"Q"+msvc.Name(),
		func(msg *nats.Msg) {
			log.Debugf("NATS %s", msg.Subject)
			ns.handleMessage(conn, msg)
		})
	if err != nil {
		panic(log.Wrapf(err, "NATS Queue Subscription failed."))
	}

	//handler.defaultReplyQ = subject+".reply"
	//return nil

	//subscribed successfully, now
	//(todo) block until subscription terminated
	//(for now block for ever)
	ch := make(chan bool)
	<- ch
} //natsServer.Run()

func (ns natsServer) handleMessage(conn *nats.Conn, msg *nats.Msg) {
	log.Debugf("Received: %s", string(msg.Data))

	//execute the operation
	responseMessage := ns.msvc.HandleJSON(operNameFromSubject(msg.Subject), msg.Data)
	jsonResponseMessage,_ := json.Marshal(responseMessage)

	if err := conn.Publish(msg.Reply, jsonResponseMessage); err != nil {
		log.Errorf("Failed to reply to \"%s\": %+v", msg.Reply, err)
	}/* else {
		log.Debugf("Replied to: \"%s\"", msg.Reply)
	}*/
}

func operNameFromSubject(subject string) string {
	parts := strings.SplitN(subject, ".", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}//operNameFromSubject()

func init() {
	//register &<struct>{} so that Validate() method will be called with pointer receiver
	//and be able to set defaults
	msvc.RegisterServer("nats", &natsServer{})
}
