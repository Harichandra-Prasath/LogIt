package LogIt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	DATE_FLAG = 1 << iota
	TIME_FLAG
	STD_FLAG = DATE_FLAG | TIME_FLAG
)

const SPACING = "  "

type Handler interface {

	// Takes the log record and prepare it for the writing
	handle(Record) error

	// Writes to the defined writer
	write(*bytes.Buffer, io.Writer) error
}

// Defualt handler which writes to os.stdout for All levels except error
type DefaultHandler struct {
	Lock sync.Mutex
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
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

// Responsible for preparing the buffer that has to be written
func (h *DefaultHandler) handle(r Record) error {

	var _buff []byte

	buff := bytes.NewBuffer(_buff)

	// Build the predefined flags
	t := time.Now()
	populateFlags(buff, t, r.Flags)

	_rc := fmt.Sprintf("%s%s%s\n", r.Level, SPACING, strings.Join(r.Message, " "))
	buff.WriteString(_rc)

	// Set the writer
	var _w io.Writer

	if r.Level == "ERROR" {
		_w = os.Stderr
	} else {
		_w = os.Stdout
	}

	err := h.write(buff, _w)
	if err != nil {
		return fmt.Errorf("on handling: %s", err)
	}

	return nil
}

// Responsible for writing the bytes to the writer
func (h *DefaultHandler) write(buff *bytes.Buffer, w io.Writer) error {

	_raw := buff.Bytes()

	h.Lock.Lock()
	_, err := w.Write(_raw)
	h.Lock.Unlock()

	if err != nil {
		return fmt.Errorf("on writing: %s", err)
	}

	return nil

}
