package fftool

import "fmt"

func errWrap(e error, msg string) error {
	return fmt.Errorf("%s:%v", msg, e)
}
