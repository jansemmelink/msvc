package main

import (
	"github.com/jansemmelink/log"
	"github.com/jansemmelink/msvc"
)

//Template create the micro-service with several operations to demonstrate how the framework is used
func Template() msvc.IMicroService {
	return msvc.New("template").
		WithOper("hello", hello{})
}

//hello operation implements msvc.IOper
type hello struct {
	msvc.Oper
	Name string `json:"name"`
}

func (h hello) Validate() error {
	if len(h.Name) == 0 {
		log.Debugf("Invalid hello: %+v", h)
		return log.Wrapf(nil, "missing name")
	}
	log.Debugf("Valid hello: %+v", h)
	return nil
}

func (h hello) Results() []msvc.IResult {
	panic("NYI")
	//return nil
}

func (h hello) Run() (interface{}, *msvc.Error) { //Run() (msvc.IResult, interface{}) {
	log.Debugf("Hello: %+v", h)
	return "Hi " + h.Name, nil
	//return nil, &msvc.Error{Type: "NYI"}
}
