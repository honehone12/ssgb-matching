package gsip

const (
	StatusOk int = 1 + iota
	StatusError
)

type GSIP struct {
	Id       string
	ClassTag int64
	Address  string
	Port     uint16
}

type GSIPResult struct {
	Status int
	Gsip   GSIP
}

type Provider interface {
	NextGsip(classTag int64) (GSIP, error)
	BackFillGsip(classTag int64) (GSIP, error)
}
