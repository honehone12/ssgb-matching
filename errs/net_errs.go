package errs

import "errors"

func ErrorDuplicatedConnection() error {
	return errors.New("duplicated connection")
}

func ErrorBrokenClock() error {
	return errors.New("clock is broken")
}
