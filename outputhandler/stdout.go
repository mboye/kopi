package outputhandler

import (
	"encoding/json"
	"os"
)

type stdoutHandler struct{ encoder *json.Encoder }

func NewStdoutHandler() OutputHandler {
	return &stdoutHandler{encoder: json.NewEncoder(os.Stdout)}
}

func (oh *stdoutHandler) Handle(obj interface{}) error {
	return oh.encoder.Encode(obj)
}
