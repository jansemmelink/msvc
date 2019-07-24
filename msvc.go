package msvc

import (
	"encoding/json"
	"reflect"

	"github.com/jansemmelink/config"
	"github.com/jansemmelink/log"
)

//IMicroService ...
type IMicroService interface {
	WithOper(name string, operTmpl IOper) IMicroService
	Serve()
	//
	Test(operName string, requestJSON string)
}

//New creates the named micro-service
func New(name string) IMicroService {
	return msvc{
		name: name,
		//default config from files in ./conf/...json|yml|properties
		configSet: config.NewSet().MustSource("files", "./conf"),
		//operations is empty until WithOper() is used
		oper: make(map[string]IOper),
	}
}

type msvc struct {
	name      string
	configSet config.ISet
	oper      map[string]IOper
}

func (msvc msvc) WithOper(name string, operTmpl IOper) IMicroService {
	if len(name) == 0 {
		panic("cannot add oper without a name")
	}
	if _, ok := msvc.oper[name]; ok {
		panic(log.Wrapf(nil, "MicroService[%s].oper[%s] already exists", msvc.name, name))
	}
	msvc.oper[name] = operTmpl
	return msvc
}

func (msvc msvc) Test(operName string, requestJSON string) {
	operTmpl, ok := msvc.oper[operName]
	if !ok {
		panic(log.Wrapf(nil, "MicroService[%s].oper[%s] does not exist", msvc.name, operName))
	}

	//create new oper instance
	operValue := reflect.New(reflect.TypeOf(operTmpl))
	operStructPtr := operValue.Interface()
	if err := json.Unmarshal([]byte(requestJSON), operStructPtr); err != nil {
		panic(log.Wrapf(err, "Failed to decode request into %T", operStructPtr))
	}

	oper := operStructPtr.(IOper)
	if err := oper.Validate(); err != nil {
		panic(log.Wrapf(err, "Invalid request"))
	}

	log.Debugf("Request is valid")
	result, response := oper.Run()

	log.Debugf("MicroService[%s].Oper[%s].Run() -> result=%v, response=(%T)%+v", msvc.name, operName, result, response, response)
	return
}

//Run starts service the micro-service on the interfaces added to the service
func (msvc msvc) Serve() {
	//see which servers are configured:
	startServers(msvc.configSet)

	panic("NYI")
}
