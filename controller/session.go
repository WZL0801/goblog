package controller

import (
	"fmt"
	"sync"
	"time"
)

// Session服务

// SessionData 表示一个具体的用户Session数据
type SessionData struct {
	ID     string
	Data   map[string]interface{}
	rwLock sync.RWMutex // 读写锁，锁的是上面的Data
	// 过期时间
}

// NewSessionData 构造函数
func NewSessionData(id string) *SessionData {
	return &SessionData{
		ID:   id,
		Data: make(map[string]interface{}, 8),
	}
}

// Mgr 是一个全局的Session 管理
type Mgr struct {
	Session map[string]SessionData
	rwLock  sync.RWMutex
}

// GetSessionData 根据传进来的SessionID找到对应的SessionData
func (m *Mgr) GetSessionData(sessionID string) (sd SessionData, err error) {
	// 取之前加锁
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	sd, ok := m.Session[sessionID]
	if !ok {
		err = fmt.Errorf("invalid session id")
		return
	}
	return
}

// CreateSession 创建一条Session记录
func (m *Mgr) CreateSession() (sd *SessionData, err error) {
	// 1. 造一个sessionID
	timeStamp := time.Now().UnixNano() // ?
	if err != nil {
		return
	}
	// 2. 造一个和它对应的SessionData
	sd = NewSessionData(string(timeStamp))
	// 3. 返回SessionData
	return
}
