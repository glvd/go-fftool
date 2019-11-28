package fftool

import "fmt"

func errWrap(e error, msg string) error {
	return fmt.Errorf("%s:%w", msg, e)
}
