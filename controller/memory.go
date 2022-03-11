package controller

import (
	"fmt"
)

// 内存版Session服务

// SessionData支持的操作

// Get 根据key获取值
func (s *SessionData) Get(key string) (value interface{}, err error) {
	// 获取读锁
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	value, ok := s.Data[key]
	if !ok {
		err = fmt.Errorf("invalid Key")
		return
	}
	return
}

// Set 根据key获取值
func (s *SessionData) Set(key string, value interface{}) {
	// 获取写锁
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.Data[key] = value
}

// Del 删除Key对应的键值对
func (s *SessionData) Del(key string) {
	// 删除key对应的键值对
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	delete(s.Data, key)
}
