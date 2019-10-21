package shardmap

import (
	"github.com/cespare/xxhash"
	"github.com/tidwall/rhh"
	"runtime"
	"sync"
)

type Map struct {
	init   sync.Once
	cap    int
	shards int
	seed   uint32
	mus    []sync.RWMutex
	maps   []*rhh.Map
}

func New(cap int) *Map {
	return &Map{cap:cap}
}

func (m *Map)initDo()  {
	m.init.Do(func() {
		m.shards = 1
		for m.shards < runtime.NumCPU()*16 {
			m.shards *= 2
		}
		scap := m.cap / m.shards
		m.mus = make([]sync.RWMutex, m.shards)
		m.maps = make([]*rhh.Map, m.shards)
		for i := 0; i < m.shards; i++ {
			m.maps[i] = rhh.New(scap)
		}
	})
}

func (m *Map)Clear()  {
	m.initDo()
	for i := 0; i < m.shards; i++ {
		m.mus[i].Lock()
		m.maps[i] = rhh.New(m.cap / m.shards)
		m.mus[i].Unlock()
	}
}

func (m *Map)choose(key string) int {
	return int(xxhash.Sum64String(key) & uint64(m.shards-1))
}

func (m *Map)Set(key string, value interface{}) (prev interface{}, replaced bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].Lock()
	prev, replaced = m.maps[shard].Set(key, value)
	m.mus[shard].Unlock()
	return
}

func (m *Map)Get(key string) (value interface{}, replaced bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].RLock()
	value, replaced = m.maps[shard].Get(key)
	m.mus[shard].RUnlock()
	return
}

func (m *Map)Delete(key string) (prev interface{}, deleted bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].Lock()
	prev, deleted = m.maps[shard].Delete(key)
	m.mus[shard].Unlock()
	return
}

func (m *Map)Len() int {
	m.initDo()
	var length int
	for i := 0; i < m.shards; i++ {
		m.mus[i].RLock()
		length += m.maps[i].Len()
		m.mus[i].RUnlock()
	}
	return length
}

func (m *Map)SetAccept(
	key string, value interface{},
	accept func(prev interface{}, replaced bool) bool,
	) (prev interface{}, replaced bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].Lock()
	defer m.mus[shard].Unlock()
	prev, replaced = m.maps[shard].Set(key, value)
	if accept != nil && !accept(prev, replaced) {
		if replaced {
			m.maps[shard].Set(key, prev)
		} else {
			m.maps[shard].Delete(key)
		}
		prev, replaced = nil, false
	}
	return
}

func (m *Map)DeleteAccept(
	key string,
	accept func(prev interface{}, deleted bool) bool,
	) (prev interface{}, deleted bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].Lock()
	defer m.mus[shard].Unlock()
	prev, deleted = m.maps[shard].Delete(key)
	if accept != nil && !accept(prev, deleted) {
		if deleted {
			m.maps[shard].Set(key, prev)
		}
		prev, deleted = nil, false
	}
	return
}

func (m *Map)Range(iter func(key string, value interface{}) bool)  {
	m.initDo()
	var done bool
	for i := 0; i < m.shards; i++ {
		m.mus[i].RLock()
		defer m.mus[i].RUnlock()
		m.maps[i].Range(func(key string, value interface{}) bool {
			if !iter(key, value) {
				done = true
				return false
			}
			return true
		})
		if done {
			break
		}
	}
}




