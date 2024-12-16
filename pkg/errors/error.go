package errors

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

const depth = 32

type stack []uintptr

type errWrapper struct {
	error
	stack *stack
}

func callers(pos int) *stack {
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[pos:n]
	return &st
}

// funcname
//func funcname(name string) string {
//	i := strings.LastIndex(name, "/")
//	name = name[i+1:]
//	n := strings.Split(name, ".")
//	return n[len(n)-1]
//}

func wrapErr(err error, caller int) *errWrapper {
	wrapper := &errWrapper{}
	wrapper.error = err
	wrapper.stack = callers(caller)

	return wrapper
}

// Wrap return an error annotating err with a stack trace
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	wrapper, ok := err.(*errWrapper)
	if !ok || wrapper == nil {
		wrapper = wrapErr(err, 1)
	}

	return wrapper
}

// New retrun new error annotated with stack trace
func New(message string) error {
	return wrapErr(errors.New(message), 1)
}

// Newf create new stack trace annotated error with formatted message
func Newf(messageFormat string, a ...interface{}) error {
	return wrapErr(fmt.Errorf(messageFormat, a...), 1)
}

// IsEqual compare errors
func IsEqual(err1, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}

	if err1 == nil || err2 == nil {
		return false
	}

	return err1.Error() == err2.Error()
}

// StackTrace return error stack trace information
func (e *errWrapper) StackTrace() errors.StackTrace {
	f := make([]errors.Frame, len(*e.stack))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame((*e.stack)[i])
	}
	return f
}
