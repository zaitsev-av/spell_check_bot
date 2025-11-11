package wer

import "fmt"

func Wer(pkg string, err error) error {
	return fmt.Errorf(" %s: %w", pkg, err)
}
