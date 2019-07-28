package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/jansemmelink/log"
	"github.com/jansemmelink/msvc"
)

//restServer implements msvc.IServer to be a HTTP REST interface for micro-services
type restServer struct {
	Address string `json:"address" doc:"HTTP Server address, e.g. localhost:12345"`

	//run-time private data:
	msvc msvc.IMicroService
}

func (rs restServer) Validate() error {
	if len(rs.Address) == 0 {
		return log.Wrapf(nil, "Missing address")
	}
	log.Debugf("Validated %T", rs)
	return nil
}

func (rs restServer) Run(msvc msvc.IMicroService) {
	rs.msvc = msvc
	http.ListenAndServe(rs.Address, rs)
}

func (rs restServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Debugf("HTTP %s %s", req.Method, req.URL)

	//read request into byte buffer
	jsonRequestData := read(req.Body)
	log.Debugf("Request: %s", string(jsonRequestData))

	// var v interface{}
	// var jsonResponseMessage []byte
	// if err := json.NewDecoder(req.Body).Decode(&v); err != nil {
	// 	jsonResponseMessage, _ = json.Marshal(msvc.ResponseMessage{
	// 		Error: &msvc.Error{
	// 			Type:        "invalidJSON in HTTP request body",
	// 			Description: err.Error(),
	// 		},
	// 	})
	// } else {
	responseMessage := rs.msvc.HandleJSON(operNameFromURL(req.URL), jsonRequestData)
	jsonResponseMessage, _ := json.Marshal(responseMessage)
	// }
	res.Write(jsonResponseMessage)
}

func operNameFromURL(url *url.URL) string {
	parts := strings.SplitN(url.Path, "/", 3)
	if len(parts) == 3 {
		//part[0] = "", part[1] = <domain> part[2] = oper
		return parts[2]
	}
	return ""
} //operNameFromURL()

func init() {
	msvc.RegisterServer("rest", restServer{})
}

func read(r io.ReadCloser) []byte {
	output := make([]byte, 0)
	buff := make([]byte, 1024)
	for {
		n, err := r.Read(buff)
		if n > 0 {
			output = append(output, buff[0:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("READ error: %v", err)
			break
		}
	}
	return output
} //read()
