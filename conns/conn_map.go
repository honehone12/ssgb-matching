package conns

import (
	"ssgb-matching/errs"
	"sync"
)

type ConnMap struct {
	count int64
	inner sync.Map
}

func NewConnMap() *ConnMap {
	return &ConnMap{
		inner: sync.Map{},
	}
}

func (m *ConnMap) Count() int64 {
	return m.count
}

func (m *ConnMap) Set(id string, conn Conn) {
	_, exists := m.inner.LoadOrStore(id, conn)
	if !exists {
		m.count++
	} else {
		_, exists = m.inner.Swap(id, conn)
		if !exists {
			panic("item should exists here")
		}
	}
}

func (m *ConnMap) Remove(id string) {
	if _, exists := m.inner.LoadAndDelete(id); exists {
		m.count--
	}
}

func (m *ConnMap) Get(id string) (Conn, error) {
	i, ok := m.inner.Load(id)
	if ok {
		c, ok := i.(Conn)
		if ok {
			return c, nil
		} else {
			return Conn{}, errs.ErrorCastFail("conn")
		}
	} else {
		return Conn{}, errs.ErrorNoSuchItem()
	}
}
