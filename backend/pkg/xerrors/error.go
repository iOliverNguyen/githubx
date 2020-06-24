package xerrors

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type M map[string]interface{}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type listErrors struct {
	Msg    string
	Errors []error
}

func Errors(msg string, errs []error) error {
	return listErrors{Msg: msg, Errors: errs}
}

func (es listErrors) Error() string {
	return fmt.Sprint(es)
}

func (es listErrors) Format(st fmt.State, c rune) {
	if es.Msg == "" && len(es.Errors) == 0 {
		_, _ = st.Write([]byte("<nil>"))
		return
	}

	width, ok := st.Width()
	if !ok {
		width = 8
	}

	verbose := st.Flag('#') || st.Flag('+')
	var b bytes.Buffer
	if es.Msg != "" {
		b.WriteString(es.Msg)
		if len(es.Errors) == 0 {
			return
		}
		if verbose {
			b.WriteString(":\n")
		} else {
			b.WriteString(": ")
		}
	}
	for i, e := range es.Errors {
		if verbose {
			for j := 0; j < width; j++ {
				b.WriteByte(' ')
			}
		}
		b.WriteString(e.Error())
		if i > 0 {
			if verbose {
				b.WriteString("\n")
			} else {
				b.WriteString("; ")
			}
		}
	}
	_, _ = st.Write(b.Bytes())
}

type stacker interface {
	StackTrace() errors.StackTrace
}

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg }
func (w *withMessage) Cause() error  { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func Errorf(err error, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	if err != nil {
		if _, ok := err.(stacker); !ok {
			err = errors.WithStack(err)
		}
		return &withMessage{
			cause: err,
			msg:   msg,
		}
	}
	return errors.New(msg)
}
