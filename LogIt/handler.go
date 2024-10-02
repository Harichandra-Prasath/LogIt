package LogIt

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

const (
	DATE_FLAG = 1 << iota
	TIME_FLAG
)

const SPACING = "  "

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

func populateFlags(buff *bytes.Buffer, t time.Time, flags int) {

	if flags&(DATE_FLAG|TIME_FLAG) != 0 {
		if flags&(DATE_FLAG) != 0 {
			y, m, d := t.Date()
			_s := fmt.Sprintf("%d/%d/%d%s", d, int(m), y, SPACING)
			buff.WriteString(_s)
		}
		if flags&(TIME_FLAG) != 0 {
			h, m, s := t.Clock()
			_s := fmt.Sprintf("%d:%d:%d%s", h, m, s, SPACING)
			buff.WriteString(_s)
		}
	}
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

	// Build the predefined flags
	t := time.Now()
	populateFlags(buff, t, r.Flags)

	_rc := fmt.Sprintf("%s%s%s\n", r.Level, SPACING, strings.Join(r.Message, " "))
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
