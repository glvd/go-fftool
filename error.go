package fftool

import (
	"fmt"
	"github.com/goextension/log"
)

// Error ...
type Error interface {
	Err() error
}

func errWrap(e error, msg string) error {
	return fmt.Errorf("%s:%w", msg, e)
}

// LogError ...
func LogError(err Error) bool {
	if err.Err() != nil {
		log.Error(err.Err())
		return true
	}
	return false
}
