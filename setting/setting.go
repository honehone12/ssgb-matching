package setting

import (
	"encoding/json"
	"io"
	"os"
)

type Setting struct {
	ServiceName     string
	ServiceVersion  string
	ServiceListenAt string

	LogLevel int

	Classes int

	QInitialCapacity   int
	RollingIntervalMil int
}

func LoadSetting(fileName string) (Setting, error) {
	s := Setting{}
	f, err := os.Open(fileName)
	if err != nil {
		return s, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(b, &s)
	return s, err
}
