package outputhandler

import (
	"encoding/json"
	"log"
	"os"
)

type stdoutHandler struct{ encoder *json.Encoder }

func NewStdoutHandler() OutputHandler {
	return &stdoutHandler{encoder: json.NewEncoder(os.Stdout)}
}

func (oh *stdoutHandler) Handle(obj interface{}) (err error) {
	err = oh.encoder.Encode(obj)
	if err != nil {
		log.Fatalf("failed to encode: %s", err.Error())
	}
	return
}
