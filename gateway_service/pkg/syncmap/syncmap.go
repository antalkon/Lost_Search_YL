package syncmap

import "sync"

type SyncMap struct {
	mutex *sync.Mutex
	data  map[string]chan string
}

func NewSyncMap() *SyncMap {
	return &SyncMap{mutex: &sync.Mutex{}, data: make(map[string]chan string)}
}

func (sm *SyncMap) Write(key string, value chan string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.data[key] = value
}

func (sm *SyncMap) Read(key string) (chan string, bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	value, ok := sm.data[key]
	if !ok {
		return nil, false
	}
	return value, true
}

func (sm *SyncMap) Delete(key string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	delete(sm.data, key)
}
