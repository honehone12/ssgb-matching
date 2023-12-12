package errs

import (
	"errors"
	"fmt"
	"strings"
)

func ErrorCastFail(name string) error {
	return fmt.Errorf("failed to cast %s", name)
}

func IsErrorCastFail(err error) bool {
	return strings.HasPrefix(err.Error(), "failed to cast")
}

func ErrorIndexOutOfRange() error {
	return errors.New("index out of range")
}

func IsErrorIndexOutOfRange(err error) bool {
	return err.Error() == "index out of range"
}
