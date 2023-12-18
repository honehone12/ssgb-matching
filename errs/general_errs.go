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
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "failed to cast")
}

func ErrorIndexOutOfRange() error {
	return errors.New("index out of range")
}

func IsErrorIndexOutOfRange(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "index out of range"
}

func ErrorNotimplemented() error {
	return errors.New("not implemented")
}

func ErrorNoSuchItem() error {
	return errors.New("no such item")
}

func ErrorItemAlreadyExists() error {
	return errors.New("item already exists")
}
