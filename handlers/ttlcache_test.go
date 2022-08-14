package handlers

import (
	"fmt"
	"testing"
	"time"
)

type testCache struct {
	t *testing.T
}

func (t testCache) Println(v ...interface{}) {
	t.t.Log(v...)
}

func (t testCache) Printf(format string, v ...interface{}) {
	t.t.Logf(format, v...)
}

func TestCreatePutToCache(t *testing.T) {

	var CacheSystem = New(1000, 1000)

	session := "bot-" + time.Now().Format("2017-09-07 17:06:04.000000")
	CacheSystem.Put(814562, sessionData{
		messageId: 0,
		sessionId: session,
		botStatus: true,
		auth:      false,
	})

	val, ok := CacheSystem.Get(814562)

	if !ok || val.messageId != 0 || val.sessionId != session {
		t.Fail()
	}

}

func TestRewriteKey(t *testing.T) {

	var CacheSystem = New(1000, 1000)

	session := "bot-" + time.Now().Format("2017-09-07 17:06:04.000000")
	CacheSystem.Put(814562, sessionData{
		messageId: 0,
		sessionId: session,
		botStatus: true,
		auth:      false,
	})

	val, ok := CacheSystem.Get(814562)

	fmt.Println("старая сессия:" + val.sessionId)

	newSession := "bot-" + time.Now().Format("2017-09-07 17:06:04.000000")

	CacheSystem.Put(814562, sessionData{
		messageId: 1,
		sessionId: newSession,
		botStatus: true,
		auth:      false,
	})

	newVal, newOk := CacheSystem.Get(814562)

	fmt.Println("новая сессия:" + newVal.sessionId)

	if val.messageId == newVal.messageId || val.sessionId == newVal.sessionId || ok != newOk {

		fmt.Println("fdsdsfdsfdfsdsfsdf")
		t.Fail()
	}

}

func TestTTl(t *testing.T) {

	var CacheSystem = New(1000, 2)

	session := "bot-" + time.Now().Format("2017-09-07 17:06:04.000000")
	CacheSystem.Put(814562, sessionData{
		messageId: 0,
		sessionId: session,
		botStatus: true,
		auth:      false,
	})

	_, ok := CacheSystem.Get(814562)

	if !ok {
		t.Fail()
	}

	time.Sleep(4 * time.Second)
	_, newOk := CacheSystem.Get(814562)

	if newOk {
		t.Fail()
	}

}

func TestDelete(t *testing.T) {

	var CacheSystem = New(1000, 1000)

	session := "bot-" + time.Now().Format("2017-09-07 17:06:04.000000")
	CacheSystem.Put(814562, sessionData{
		messageId: 0,
		sessionId: session,
		botStatus: true,
		auth:      false,
	})

	_, ok := CacheSystem.Get(814562)

	if !ok {
		t.Fail()
	}

	CacheSystem.Delete(814562)

	_, newok := CacheSystem.Get(814562)

	if newok {
		t.Fail()
	}

}
