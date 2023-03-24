package dwebsocket

import "sync"

var ClientMap = NewCliMap()

type CliMap struct {
	sync.RWMutex
	m map[int64]*Client
}

func NewCliMap() *CliMap {
	return &CliMap{
		m: make(map[int64]*Client),
	}
}

func (m *CliMap) Get(k int64) (*Client, bool) {
	m.RLock()
	defer m.RUnlock()
	v, existed := m.m[k]
	return v, existed
}

func (m *CliMap) Set(k int64, v *Client) {
	m.Lock()
	defer m.Unlock()
	m.m[k] = v
}

func (m *CliMap) Delete(k int64) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, k)
}

func (m *CliMap) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.m)
}

func (m *CliMap) Each(f func(k int64, v *Client) bool) {
	m.RLock()
	defer m.RUnlock()

	for k, v := range m.m {
		if !f(k, v) {
			return
		}
	}
}
