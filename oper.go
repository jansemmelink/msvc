package msvc

import (
	"fmt"
	"strings"
)

//IOper is one operation in the micro-service
type IOper interface {
	//Results return a list of results that this operation may return
	Results() []IResult

	//Validate the operation request before it is called
	Validate() error

	//Run the operation to return the (optional) response data or an error
	Run() (interface{}, *Error)

	//ErrorMessage can be used on a request to respond with an error
	ErrorMessage(errorType string, err error) ResponseMessage
}

//Oper ...
type Oper struct{}

//ErrorMessage ...
func (oper Oper) ErrorMessage(errorType string, err error) ResponseMessage {
	errText := strings.TrimPrefix(fmt.Sprintf("%s", err), "because ")
	return ResponseMessage{
		Error: &Error{
			Type:        errorType,
			Description: errText,
		},
	}
}
