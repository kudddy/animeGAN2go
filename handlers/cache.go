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

// `SessionData
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
	// TODO i am not sure what is optimal method, but it is work
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
		it.value.botStatus = true
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

//var CacheSystem = New(1000, 1000)

//var CacheSystemUser = New(1000, 1000)
//var CacheSystemOperator = New(1000, 1000)

type CachePolicy struct {
	cache map[string]*TTLMap
}

func (l *CachePolicy) AddData(projectId string, Cache *TTLMap) {
	l.cache[projectId] = Cache
}

func (l *CachePolicy) GetData(projectIs string) (*TTLMap, bool) {
	d, ok := l.cache[projectIs]
	return d, ok
}
func (l *CachePolicy) Init() {
	l.cache = map[string]*TTLMap{}
}

// CACHE FOR BOT PARAMS BY PROJECT
type botsInfo struct {
	bot      string
	operator string
	webhook  string
}

var BotsInfo = botsInfo{
	"5590021672:AAFi97mI_4hcJK2YjPe7xYkO5KjtyncnzGc",
	"5409018161:AAE7hHy1C3cbmiNAvTTpT59AVYNH1_nFAVQ",
	"https://smartapp-code.sberdevices.ru/chatadapter/chatapi/webhook/sber_nlp2/tFPVKJfU:a43f2e017c60069c99b387ee5a9839eebe490999",
}

type botsParams struct {
	data map[string]botsInfo
}

func (l *botsParams) GetData(projectIs string) (botsInfo, bool) {
	d, ok := l.data[projectIs]
	return d, ok
}

func (l *botsParams) AddData(projectId string, botData botsInfo) {
	l.data[projectId] = botData
}

func (l *botsParams) Init() {
	l.data = map[string]botsInfo{}
}

func Init() botsParams {
	var m botsParams
	m.Init()
	// test data, in future we must get it from base
	m.AddData("7d216f7c-cfee-4b76-a550-1c66a93848c9", BotsInfo)
	return m
}

func InitCache() CachePolicy {
	var cache CachePolicy
	cache.Init()
	for _, projectId := range AuthTokens {
		cache.AddData(projectId, New(1000, 1000))
	}
	return cache
}

var BotsParams = Init()
var CacheUser = InitCache()
var CacheOperator = InitCache()

// AuthTokens temporary cache for auth, in future this data we will get from database
var AuthTokens = []string{"7d216f7c-cfee-4b76-a550-1c66a93848c9", "b1630dbc-51a4-4462-81c8-5233d2a92081"}
