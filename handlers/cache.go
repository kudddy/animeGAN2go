package handlers

import (
	"log"
	"sync"
	"time"
)

type Dictionary map[string]struct {
	MessageId int
	sessionId string
	botStatus bool
}

var BotsInfo = map[string]string{
	"bot":      "***",
	"operator": "***",
}

// CACHE SYSTEM
type sessionData struct {
	sessionId       string
	botStatus       bool
	messageId       int
	companionUserId int
}

type item struct {
	value      sessionData
	lastAccess int64
}

type TTLMap struct {
	m map[int]*item
	l sync.Mutex
}

func New(ln int, maxTTL int) (m *TTLMap) {
	m = &TTLMap{m: make(map[int]*item, ln)}
	go func() {
		for now := range time.Tick(time.Second) {
			m.l.Lock()
			for k, v := range m.m {
				if now.Unix()-v.lastAccess > int64(maxTTL) {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()
	return
}

func (m *TTLMap) Len() int {
	return len(m.m)
}

func (m *TTLMap) IterMid(k int) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		it.value.messageId = +1
		m.m[k] = it

	}
	m.l.Unlock()

	return
}

func (m *TTLMap) ChangeBotStatus(k int) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		it.value.botStatus = false
		m.m[k] = it
	}
	m.l.Unlock()

	return
}

func (m *TTLMap) Put(k int, v sessionData) {
	m.l.Lock()

	log.Printf("save session parametrs for id is %d, companion id is %d . log from PUT", k, v.companionUserId)
	it, _ := m.m[k]

	it = &item{value: v}
	m.m[k] = it

	it.lastAccess = time.Now().Unix()
	m.l.Unlock()
}

func (m *TTLMap) Get(k int) (v sessionData, found bool) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		v = it.value
		found = ok
		it.lastAccess = time.Now().Unix()
	}
	m.l.Unlock()
	return

}

func (m *TTLMap) Delete(k int) {
	m.l.Lock()
	delete(m.m, k)
	m.l.Unlock()
}

var CacheSystem = New(1000, 30)
