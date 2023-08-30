package service

import (
	"github.com/gomodule/redigo/redis"
)

type RedisService interface {
	GetByKey(id string) (string, error)
	GetAll() ([]string, error)
	GetAllKeys() ([]string, error)
	GetPrefix() string
	GetById(id string) (string, error)
}

type JsonService struct {
	prefix    string
	redisPool *redis.Pool
}

func NewJsonService(prefix string, redisPool *redis.Pool) *JsonService {
	return &JsonService{prefix, redisPool}
}

func (s *JsonService) GetAll() ([]string, error) {
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

func (s *JsonService) GetPrefix() string {
	return s.prefix
}

func (s *JsonService) GetById(id string) (string, error) {
	return s.GetByKey(s.prefix + id)
}

func (s *JsonService) GetByKey(key string) (string, error) {
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

func (s *JsonService) GetAllKeys() ([]string, error) {
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
