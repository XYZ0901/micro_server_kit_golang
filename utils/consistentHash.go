package utils

import (
	consul "github.com/hashicorp/consul/api"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Uint32Slice []uint32

func (s Uint32Slice) Len() int {
	return len(s)
}

func (s Uint32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Uint32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	key      Uint32Slice
	hashMap  map[uint32]string
	sync.RWMutex
}

func NewMap(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[uint32]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) IsEmpty() bool {
	return len(m.key) == 0
}

func (m *Map) Add(keys ...string) {
	m.Lock()
	defer m.Unlock()
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := m.hash([]byte(strconv.Itoa(i) + key))
			m.key = append(m.key, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Sort(m.key)
}

func (m *Map) Delete(keys ...string) {
	m.Lock()
	defer m.Unlock()
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := m.hash([]byte(strconv.Itoa(i) + key))
			if _, ok := m.hashMap[hash]; !ok {
				return
			}
			delete(m.hashMap, hash)
			idx := sort.Search(m.key.Len(), func(i int) bool {
				return m.key[i] == hash
			})
			m.key = append(m.key[:idx], m.key[idx:]...)
		}
	}
}

func (m *Map) Get(key string) string {
	m.RLock()
	defer m.RUnlock()
	if m.IsEmpty() {
		return ""
	}
	hash := m.hash([]byte(key))
	idx := sort.Search(m.key.Len(), func(i int) bool {
		return m.key[i] >= hash
	})
	if idx == len(m.key) {
		idx = 0
	}
	return m.hashMap[m.key[idx]]
}

func ChooseServerFromMap(as map[string]*consul.AgentService) consul.AgentService {
	m := NewMap(10, nil)
	keys := []string{}
	for k, _ := range as {
		keys = append(keys, k)
	}
	m.Add(keys...)
	key := m.Get(time.Now().String())
	return *as[key]
}
