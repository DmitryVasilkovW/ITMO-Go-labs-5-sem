//go:build !solution

package keylock

import (
	"sort"
	"sync"
)

type bc = chan bool

type KeyLock struct {
	mutex *sync.Mutex
	m     map[string]bc
}

func New() *KeyLock {
	return &KeyLock{mutex: &sync.Mutex{}, m: make(map[string]bc)}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	key := prepare(keys)
	var r bool
	unlock = func() {
		l.releaseKeys(key)
	}

	for _, k := range key {
		if l.tryLockKey(k, cancel) {
			r = true
			unlock()
			return r, unlock
		}
	}

	return r, unlock
}

func (l *KeyLock) tryLockKey(key string, cancel <-chan struct{}) bool {
	l.mutex.Lock()
	value := l.getOrCreateChannel(key)
	l.mutex.Unlock()

	return l.waitForChannel(value, cancel)
}

func prepare(keys []string) []string {
	key := make([]string, len(keys))
	copy(key, keys)
	sort.Strings(key)
	return key
}

func (l *KeyLock) releaseKeys(keys []string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for _, k := range keys {
		l.releaseChannel(k)
	}
}

func (l *KeyLock) releaseChannel(key string) {
	select {
	case l.m[key] <- true:
	default:
	}
}

func (l *KeyLock) getOrCreateChannel(key string) bc {
	if c, exists := l.m[key]; exists {
		return c
	}

	c := make(bc, 1)
	c <- true
	l.m[key] = c
	return c
}

func (l *KeyLock) waitForChannel(ch bc, cancel <-chan struct{}) bool {
	select {
	case <-cancel:
		return true
	case <-ch:
		return false
	}
}
