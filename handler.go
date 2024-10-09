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

var COLOR_MAP = map[int]string{
	DATE_FLAG: "\033[38;5;209m",
	TIME_FLAG: "\033[38;5;154m",
	0:         "\033[38;5;104m",
	-1:        "\033[39m",
}

type Handler interface {

	// Takes the Log record and prepare it for the writing
	handle(Record) error

	// Writes to the defined writer
	write(*bytes.Buffer, bool) error
}

// Basic Text Handler for logger
type TextHandler struct {
	Lock sync.Mutex

	// writer to write logs other than level ERROR
	Out io.Writer

	// writer to write logs for level ERROR
	Err io.Writer
}

// Returns a new TextHandler for given writers
func NewTextHandler(out io.Writer, err io.Writer) *TextHandler {
	return &TextHandler{
		Out: out,
		Err: err,
	}
}

// Populates requested flags before writing the main content of the logs
func populateFlags(buff *bytes.Buffer, t time.Time, flags int, spacing string, colorfull bool) {

	if flags&(DATE_FLAG|TIME_FLAG) != 0 {
		if flags&(DATE_FLAG) != 0 {
			y, m, d := t.Date()

			if colorfull {
				buff.WriteString(COLOR_MAP[DATE_FLAG])

			}

			_s := fmt.Sprintf("%d/%d/%d%s", d, int(m), y, spacing)
			buff.WriteString(_s)
		}
		if flags&(TIME_FLAG) != 0 {

			h, m, s := t.Clock()

			if colorfull {
				buff.WriteString(COLOR_MAP[TIME_FLAG])
			}

			_s := fmt.Sprintf("%d:%d:%d%s", h, m, s, spacing)
			buff.WriteString(_s)
		}
	}

}

// Responsible for preparing the buffer that has to be written
func (h *TextHandler) handle(r Record) error {

	var _buff []byte
	var _err bool = false

	if r.Level == "ERROR" {
		_err = true
	}

	buff := bytes.NewBuffer(_buff)

	_spacing := strings.Repeat(" ", r.Options.Spacing)

	// Build the predefined flags
	t := time.Now()
	populateFlags(buff, t, r.Options.Flags, _spacing, r.Options.Colorfull)

	var level string

	if r.Options.Colorfull {
		_lc := COLOR_MAP[0]
		_re := COLOR_MAP[-1]
		level = fmt.Sprintf("%s%s%s", _lc, r.Level, _re)
	} else {
		level = r.Level
	}

	_rc := fmt.Sprintf("%s%s%s\n", level, _spacing, strings.Join(r.Message, " "))

	buff.WriteString(_rc)

	err := h.write(buff, _err)
	if err != nil {
		return fmt.Errorf("on handling: %s", err)
	}

	return nil
}

// Responsible for writing the bytes to the writer
func (h *TextHandler) write(buff *bytes.Buffer, IsErr bool) error {

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
