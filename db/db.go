package db

import "sync"

type Db struct {
	Name string
	m    map[string]int
	mu   *sync.RWMutex
}

func EmptyDb(name string) *Db {
	return &Db{name, make(map[string]int), new(sync.RWMutex)}
}

func NewDb(name string, data map[string]int) *Db {
	return &Db{name, data, new(sync.RWMutex)}
}

func (db *Db) Get(key string) int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.m[key]
}

func (db *Db) Add(key string, val int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.m[key] = val
}

func (db *Db) Reset() {
	db.m = make(map[string]int)
}
