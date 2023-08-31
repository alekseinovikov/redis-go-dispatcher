package service

import (
	"github.com/gomodule/redigo/redis"
)

type JsonServiceImpl struct {
	prefix    string
	redisPool *redis.Pool
}

func NewJsonService(prefix string, redisPool *redis.Pool) *JsonServiceImpl {
	return &JsonServiceImpl{prefix, redisPool}
}

func (s *JsonServiceImpl) GetAll() ([]string, error) {
	conn := s.redisPool.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	keys, err := s.GetAllKeys()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(keys))
	for _, key := range keys {
		data, err := redis.String(conn.Do("GET", key))
		if err != nil {
			return nil, err
		}

		result = append(result, data)
	}

	return result, nil
}

func (s *JsonServiceImpl) GetPrefix() string {
	return s.prefix
}

func (s *JsonServiceImpl) GetById(id string) (string, error) {
	return s.GetByKey(s.prefix + id)
}

func (s *JsonServiceImpl) GetByKey(key string) (string, error) {
	conn := s.redisPool.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	result, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", nil
	}

	data, err := redis.String(result, err)
	if err != nil {
		return "", err
	}

	return data, nil
}

func (s *JsonServiceImpl) GetAllKeys() ([]string, error) {
	conn := s.redisPool.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	keys, err := redis.Strings(conn.Do("KEYS", s.prefix+"*"))
	if err != nil {
		return nil, err
	}

	return keys, nil
}
