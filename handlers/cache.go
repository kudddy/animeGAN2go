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

// struct for bot param
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
	auth            bool
	busy            bool
	newSession      bool
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

func (m *TTLMap) ChangeAuthStatus(k int) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		it.value.auth = true
		m.m[k] = it
	}
	m.l.Unlock()

	return
}

func (m *TTLMap) ChangeBusyStatus(k int) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		it.value.busy = false
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

func (m *TTLMap) GetRandomAuthOperators() []int {
	var s []int
	for key, value := range m.m {
		if value.value.auth && !value.value.busy {
			s = append(s, key)
		}
	}
	return s
}

func (m *TTLMap) ChangeSessionStatus(k int) {
	m.l.Lock()
	if it, ok := m.m[k]; ok {
		it.value.newSession = false
		m.m[k] = it
	}
	m.l.Unlock()
}

func (m *TTLMap) Delete(k int) {
	m.l.Lock()
	delete(m.m, k)
	m.l.Unlock()
}

var CacheSystem = New(1000, 1000)

// CACHE FOR BOT PARAMS BY PROJECT ID
type botsParams struct {
	data map[string]map[string]string
}

func (l *botsParams) GetData(projectIs string) (map[string]string, bool) {
	d, ok := l.data[projectIs]
	return d, ok
}

func (l *botsParams) AddData(projectId string, botData map[string]string) {
	l.data[projectId] = botData
}

func (l *botsParams) Init() {
	l.data = map[string]map[string]string{}
}

func Init() botsParams {
	var m botsParams
	m.Init()
	// test data, in future we must get it from base
	m.AddData("7d216f7c-cfee-4b76-a550-1c66a93848c9", BotsInfo)
	return m
}

var BotsParams = Init()
