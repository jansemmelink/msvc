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

//hello implements msvc.IOper
type hello struct {
	Name string `json:"name"`
}

func (h hello) Validate() error {
	if len(h.Name) == 0 {
		return log.Wrapf(nil, "missing name")
	}
	return nil
}

func (h hello) Results() []msvc.IResult {
	panic("NYI")
	//return nil
}

func (h hello) Run() (msvc.IResult, interface{}) {
	log.Debugf("Hello: %+v", h)
	return nil, "Hi " + h.Name
	//panic("NYI")
	//return nil
}
