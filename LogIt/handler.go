package LogIt

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type Handler interface {

	// Takes the log record and prepare it for the writing
	Handle(Record) error

	// Writes to the defined writer
	Write(*bytes.Buffer) error
}

type DefaultHandler struct {
	Lock   sync.Mutex
	Writer io.Writer
}

func NewDefaultHandler(w io.Writer) *DefaultHandler {

	return &DefaultHandler{
		Writer: w,
	}

}

// Responsible for preparing the buffer that has to be written
func (h *DefaultHandler) Handle(r Record) error {

	var _buff []byte

	buff := bytes.NewBuffer(_buff)
	_rc := fmt.Sprintf("%s\t%s\n", r.Level, r.Message)
	buff.WriteString(_rc)

	err := h.Write(buff)
	if err != nil {
		return fmt.Errorf("on handling: %s", err)
	}

	return nil
}

// Responsible for writing the bytes to the writer
func (h *DefaultHandler) Write(buff *bytes.Buffer) error {

	_raw := buff.Bytes()

	h.Lock.Lock()
	_, err := h.Writer.Write(_raw)
	h.Lock.Unlock()

	if err != nil {
		return fmt.Errorf("on writing: %s", err)
	}

	return nil

}
