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

	STD_FLAG = DATE_FLAG | TIME_FLAG
)

type Handler interface {

	// Takes the log record and prepare it for the writing
	handle(Record) error

	// Writes to the defined writer
	write(*bytes.Buffer, bool) error
}

type LogHandler struct {
	Lock sync.Mutex
	Out  io.Writer
	Err  io.Writer
}

func NewLogHandler(out io.Writer, err io.Writer) *LogHandler {
	return &LogHandler{
		Out: out,
		Err: err,
	}
}

func populateFlags(buff *bytes.Buffer, t time.Time, flags int, spacing string) {

	if flags&(DATE_FLAG|TIME_FLAG) != 0 {
		if flags&(DATE_FLAG) != 0 {
			y, m, d := t.Date()
			_s := fmt.Sprintf("%d/%d/%d%s", d, int(m), y, spacing)
			buff.WriteString(_s)
		}
		if flags&(TIME_FLAG) != 0 {
			h, m, s := t.Clock()
			_s := fmt.Sprintf("%d:%d:%d%s", h, m, s, spacing)
			buff.WriteString(_s)
		}
	}
}

// Responsible for preparing the buffer that has to be written
func (h *LogHandler) handle(r Record) error {

	var _buff []byte
	var _err bool = false

	if r.Level == "ERROR" {
		_err = true
	}

	buff := bytes.NewBuffer(_buff)

	_spacing := strings.Repeat(" ", r.Options.Spacing)

	// Build the predefined flags
	t := time.Now()
	populateFlags(buff, t, r.Options.Flags, _spacing)

	_rc := fmt.Sprintf("%s%s%s\n", r.Level, _spacing, strings.Join(r.Message, " "))
	buff.WriteString(_rc)

	err := h.write(buff, _err)
	if err != nil {
		return fmt.Errorf("on handling: %s", err)
	}

	return nil
}

// Responsible for writing the bytes to the writer
func (h *LogHandler) write(buff *bytes.Buffer, IsErr bool) error {

	_raw := buff.Bytes()

	var _writer io.Writer

	if IsErr {
		_writer = h.Err
	} else {
		_writer = h.Out
	}

	h.Lock.Lock()
	_, err := _writer.Write(_raw)
	if err != nil {
		return fmt.Errorf("on writing: %s", err)
	}
	h.Lock.Unlock()

	return nil

}
