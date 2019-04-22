package outputhandler

type OutputHandler interface {
	Handle(obj interface{}) error
}

var Stdout = NewStdoutHandler()
