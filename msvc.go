package msvc

import (
	"encoding/json"
	"reflect"
	"sync"

	"github.com/jansemmelink/config"
	"github.com/jansemmelink/log"
)

//IMicroService ...
type IMicroService interface {
	Name() string
	WithOper(name string, operTmpl IOper) IMicroService
	Serve()
	HandleJSON(operName string, jsonRequestMessage []byte) ResponseMessage
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
		operTmpl: make(map[string]IOper),
	}
}

type msvc struct {
	name      string
	configSet config.ISet
	operTmpl  map[string]IOper
}

func (msvc msvc) Name() string {
	return msvc.name
}

func (msvc msvc) WithOper(name string, operTmpl IOper) IMicroService {
	if len(name) == 0 {
		panic("cannot add oper without a name")
	}
	if _, ok := msvc.operTmpl[name]; ok {
		panic(log.Wrapf(nil, "MicroService[%s].oper[%s] already exists", msvc.name, name))
	}
	msvc.operTmpl[name] = operTmpl
	return msvc
}

func (msvc msvc) Test(operName string, requestJSON string) {
	operTmpl, ok := msvc.operTmpl[operName]
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

//Serve the micro-service on all the configured server interfaces
func (msvc msvc) Serve() {
	//start all the configured servers
	wg := sync.WaitGroup{}
	startConfiguredServers(&wg, msvc.configSet, msvc)

	//wait for all servers to terminate
	wg.Wait()
	log.Debugf("All servers terminated.")
}

//HandleJSON is called by all the IServer implementations when they received a JSON message
func (msvc msvc) HandleJSON(operName string, jsonRequestMessage []byte) ResponseMessage {
	operTmpl, ok := msvc.operTmpl[operName]
	if !ok {
		return ResponseMessage{
			Header:   nil,
			Request:  nil,
			Error:    &Error{Type: "unknownOper"},
			Response: nil}
	}

	//startTime := time.Now()
	// monitor.GaugeInc("concurrent_transactions", "")
	// defer monitor.GaugeDec("concurrent_transactions", "")

	//decode only {"header":{...}}, ignoring the rest of the request message
	var requestMessage RequestMessage
	if err := fromJSON(requestMessage.RequestMessageOnlyHeader, jsonRequestMessage); err != nil {
		return ResponseMessage{
			Error: &Error{
				Type:        "decodeJSONRequestHeader",
				Description: log.Wrapf(err, "Failed to decode request header").Error(),
			},
		}
	}

	log.Debugf(".Header: %+v", requestMessage)

	if /*requestTimestamp*/ _ /*maxDur*/, _, err := requestMessage.Validate(operName); err != nil {
		log.Debugf("Invalid request message")
		return ResponseMessage{
			Error: &Error{
				Type:        "invalidRequestHeader",
				Description: log.Wrapf(err, "Invalid request header").Error(),
			},
		}
	}
	log.Debugf("Valid request message: %+v", requestMessage)

	//reject requests when terminating
	// if terminating {
	// 	return "", "", ProcessIsTerminating, errors.Errorf("Process is terminating")
	// }

	// var result *Result
	// var resultName string
	// var msAudit interface{}

	/*
	 * Log the transaction in prometheus and in the audit log
	 */
	// defer func() {

	// 	endTime := time.Now()

	// 	RecordTransaction(
	// 		domain,
	// 		operation,
	// 		message.Header.IntGuid,
	// 		startTime,
	// 		endTime,
	// 		result,
	// 		resultName,
	// 		msAudit)

	// }()

	/*
	 * Return if request validation failed
	 */
	// if err != nil {

	// 	err = errors.Wrapf(err,
	// 		"message validation failed")
	// 	_log.Errorf("%+v", err)

	// 	result = MakeResultFromResultCode(resultCode, err.Error())
	// 	message.Header.Result = result

	// 	ms.connector.Reply(message)
	// 	return

	// } // if message validation failed

	//include request.header.uuid in all subsequent logging
	//if absent, assign own unique id
	// _log = _log.With(
	// 	"guid",
	// 	message.Header.IntGuid)

	// Catch the panic
	// defer func() {

	// 	if r := recover(); r != nil {
	// 		err, ok := r.(error)
	// 		if ok && err != nil {
	// 			/*
	// 				Get stack trace info
	// 			*/
	// 			const size = 64 << 10 //64k
	// 			buf := make([]byte, size)
	// 			buf = buf[:runtime.Stack(buf, false)]
	// 			_log.Errorf("System panic with error: %v\n%s", err, buf)
	// 			message.Header.Result = MakeResult(-1, err.Error(), err.Error())
	// 		} else {
	// 			/*
	// 				Get stack trace info
	// 			*/
	// 			const size = 64 << 10 //64k
	// 			buf := make([]byte, size)
	// 			buf = buf[:runtime.Stack(buf, false)]
	// 			_log.Errorf("System panic:\n%s %s", err, buf)
	// 			message.Header.Result = MakeResult(-1, r.(string), r.(string))
	// 		}

	// 		json, err := message.ToJSON()

	// 		if err != nil {
	// 			_log.Errorf("%v", err)
	// 			return
	// 		}

	// 		_log.Debugf("Panic Message to send back is %s", json)

	// 		if len(message.Header.ReplyAddress) > 0 {
	// 			err = ms.connector.Reply(message)
	// 			if err != nil {
	// 				_log.Errorf("Failed to reply to %s, error is %v", message.Header.ReplyAddress, err)
	// 			}
	// 		}
	// 	}
	// }()

	//create a new copy of the operation (the request) struct
	operType := reflect.TypeOf(operTmpl)
	log.Debugf("operType = %v", operType)
	var operStructPtrValue reflect.Value
	if operType.Kind() == reflect.Ptr {
		operStructPtrValue = reflect.New(operType.Elem())
	} else {
		operStructPtrValue = reflect.New(operType)
	}

	//decode the message.request element into the operation
	requestMessage.Request = operStructPtrValue.Interface()
	if err := fromJSON(&requestMessage.RequestMessageOnlyRequest, jsonRequestMessage); err != nil {
		return ResponseMessage{
			Error: &Error{
				Type:        "decodeJSONRequestData",
				Description: log.Wrapf(err, "Failed to decode request data").Error(),
			},
		}
	}
	log.Debugf("Decoded request in message: %+v", requestMessage)

	operRequest, ok := requestMessage.Request.(IOper)
	if !ok {
		return ResponseMessage{
			Error: &Error{
				Type:        "operMissingValidator",
				Description: log.Wrapf(nil, "Internal Software Error: Validator not implemented").Error(),
			},
		}
	}
	log.Debugf("Got request: %+v", operRequest)

	if err := operRequest.Validate(); err != nil {
		return operRequest.ErrorMessage("invalidRequest", log.Wrapf(err, "Invalid Request"))
	}
	log.Debugf("Valid request: %+v", operRequest)

	operResponse, operError := operRequest.Run()
	if operError != nil {
		return ResponseMessage{
			Error: operError,
		}
	}
	return ResponseMessage{
		Error:    nil,
		Response: operResponse,
	}
} //msvc.HandleJSON()
