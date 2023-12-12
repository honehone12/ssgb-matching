package uuid

import (
	"errors"

	libuuid "github.com/google/uuid"
)

func ZeroUuid(id string) error {
	u, err := libuuid.Parse(id)
	if err != nil {
		return err
	}

	for i := 0; i < 16; i++ {
		if u[i] != 0 {
			return nil
		}
	}

	return errors.New("zero uuid")
}
