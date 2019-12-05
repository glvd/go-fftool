package fftool

import (
	"github.com/goextension/log"
	"strings"
)

type errWrap struct {
	err error
	msg string
}

// Error ...
func (e *errWrap) Error() string {
	return strings.Join([]string{e.msg, e.err.Error()}, ":")
}

// Err ...
func Err(e error, msg string) error {
	return &errWrap{
		err: e,
		msg: msg,
	}
}

// LogError ...
func LogError(err error) bool {
	if err != nil {
		log.Error(err.Error())
		return true
	}
	return false
}
