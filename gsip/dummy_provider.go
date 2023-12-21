package gsip

import (
	"fmt"
	"math"
	"math/rand"
	"ssgb-matching/errs"
)

type DummyProvider struct {
	id      string
	address string
	port    uint16
	count   uint16
}

func NewDummyProvider() *DummyProvider {
	return &DummyProvider{
		id:      "12345678-1234-1234-1234-123456123456",
		address: "://127.0.0.1",
		port:    7777,
		count:   0,
	}
}

func (p *DummyProvider) NextGsip(classTag int64) (GSIP, error) {
	n := rand.Intn(100)
	if n == 7 || n == 77 {
		return GSIP{}, fmt.Errorf("error [%d]", n)
	}

	gsip := GSIP{
		Id:       p.id,
		ClassTag: classTag,
		Address:  "new" + p.address,
		Port:     p.port + p.count,
	}
	p.count = (p.count + 1) % math.MaxUint16

	return gsip, nil
}

func (p *DummyProvider) BackFillGsip(classTag int64) (GSIP, error) {
	n := rand.Intn(100)
	if n == 7 || n == 77 {
		return GSIP{}, fmt.Errorf("error [%d]", n)
	}

	if n%2 == 0 {
		return GSIP{}, errs.ErrorNotFound()
	}

	gsip := GSIP{
		Id:       p.id,
		ClassTag: classTag,
		Address:  "backfill" + p.address,
		Port:     p.port + p.count,
	}
	p.count = (p.count + 1) % math.MaxUint16

	return gsip, nil
}
