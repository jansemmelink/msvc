package msvc

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/jansemmelink/log"
)

//fromJSON is used to decode JSON data into the specified struct
func fromJSON(output interface{}, jsonMessage []byte) error {
	decoder := json.NewDecoder(strings.NewReader(string(jsonMessage)))
	decoder.UseNumber()
	log.Debugf("Decoding into output %T ...", output)
	err := decoder.Decode(&output)
	if err != nil {
		return err
	}
	return nil
} //fromJSON()

//Header is common in all messages, but optional :-)
type Header struct {
	Timestamp string    `json:"timestamp" doc:"Timestamp when this message is sent written as ..."`
	UUID      string    `json:"uuid,omitempty" doc:"Optional UUID. If present in request it is echoed in the response."`
	Consumer  *Consumer `json:"consumer,omitempty" doc:"In request, describes the sender of the request. Echoed exactly in the response."`
	Provider  *Provider `json:"provider,omitempty" doc:"In request, describes who should provide the service. May be omitted if message was sent to service and operation name."`
}

//RequestMessage is separated into two embedded structures that allows us to decode
// them one at a time: the header first and later the request data...
type RequestMessage struct {
	RequestMessageOnlyHeader
	RequestMessageOnlyRequest
}

//RequestMessageOnlyHeader ...
type RequestMessageOnlyHeader struct {
	Header *RequestHeader `json:"header,omitempty" doc:"Header is optional in a request."`
}

//RequestMessageOnlyRequest ...
type RequestMessageOnlyRequest struct {
	Request interface{} `json:"request,omitempty" doc:"Request data may be nil if not present in the request."`
}

//RequestHeader ...
type RequestHeader struct {
	Header
	//header values used in request only:
	MaxDur      time.Duration `json:"max-duration" doc:"Indicate how long sender will wait for a response."`
	EchoRequest bool          `json:"echo-request" doc:"True if request data must be echoed in the response message."`
}

//Validate the request message header ...
//return:
//	timestamp
//	max-duration
//	error if not valid
func (requestMessage RequestMessage) Validate(operName string) (time.Time, time.Duration, error) {
	if requestMessage.Header == nil {
		return time.Now(), 0, nil //header is optional, so proceed if not present
	}
	h := requestMessage.Header

	// if len(h.Timestamp) == len("2017-06-07 11:37:58") {
	// 	message.Header.Timestamp = message.Header.Timestamp + ".000"
	// }
	// timestamp, err := time.ParseInLocation(
	// 	"2006-01-02 15:04:05.000",
	// 	message.Header.Timestamp,
	// 	time.Local)
	timestamp, err := time.Parse("2006-01-02 15:04:05.000+07:00", h.Timestamp)
	if err != nil {
		//try without milliseconds
		timestamp, err = time.Parse("2006-01-02 15:04:05+07:00", h.Timestamp)
		if err != nil {
			//try local time with milliseconds
			timestamp, err = time.ParseInLocation("2006-01-02 15:04:05.00", h.Timestamp, time.Local)
			if err != nil {
				//try local time without milliseconds
				timestamp, err = time.ParseInLocation("2006-01-02 15:04:05", h.Timestamp, time.Local)
				if err != nil {
					return time.Now(), 0, log.Wrapf(nil, "Invalid timestamp. Expecting %s", time.Now().Format("2006-01-02 15:04:05.000+07:00"))
				}
			}
		}
	}

	//h.UUID and h.Consumer needs no validation - echo whatever we got in the response

	//h.Provider in request indicates who 'should' handle this, but we already know the operation name
	//so we do not validate it yet...

	//h.MaxDur ...
	//h.EchoRequest ...

	//message expired if ttl > 0 and header.ts+header.ttl < now
	if h.MaxDur > 0 && time.Now().After(timestamp.Add(h.MaxDur)) {
		return timestamp, h.MaxDur, log.Wrapf(nil, "timestamp:\"%s\" + max-dur:%v has expired", h.Timestamp, h.MaxDur)
	}
	return timestamp, h.MaxDur, nil
} //RequestMessage.Validate()

//ResponseMessage ...
type ResponseMessage struct {
	Header  *ResponseHeader `json:"header,omitempty"`
	Request interface{}     `json:"request,omitempty" doc:"Request data is only present here if specified echo-request:true in the request message."`

	//followed by either error or response:
	Error    *Error      `json:"error,omitempty" doc:"Error only if service failed, then there will be no response."`
	Response interface{} `json:"response,omitempty" doc:"Response data if service succeeded and provided response data."`
}

//ResponseHeader ...
type ResponseHeader struct {
	Header
	//header values used in response only:
	Dur time.Duration `json:"duration" doc:"Duration between timestamp in request and response"`
}

//Consumer ...
type Consumer struct {
	Name string `json:"name,omitempty" doc:"Name to identify the consumer"`
	TID  string `json:"tid,omitempty" doc:"Optional Transaction ID, significant only in the context of the consumer."`
	SID  string `json:"sid,omitempty" doc:"Optional Session ID, significant only in the context of the consumer."`
}

//Provider ...
type Provider struct {
	Name string `json:"name,omitempty" doc:"Name to identify the provider"`
	TID  string `json:"tid,omitempty" doc:"Optional Transaction ID, significant only in the context of the provider."`
	SID  string `json:"sid,omitempty" doc:"Optional Session ID, significant only in the context of the provider."`
}

//Error ...
type Error struct {
	Type        string `json:"type,omitempty" doc:"Type of error is a name to identify the error in a lookup table."`
	Description string `json:"description,omitempty" doc:"Free format text to further explain the error."`
}
