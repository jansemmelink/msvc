package msvc

//IOper is one operation in the micro-service
type IOper interface {
	//Results return a list of results that this operation may return
	Results() []IResult

	//Validate the operation request before it is called
	Validate() error

	//Run the operation to return the result and response data
	Run() (IResult, interface{})
}
