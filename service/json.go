package service

import (
	"github.com/gomodule/redigo/redis"
	"strings"
)

type JsonService struct {
	prefix    string
	redisPool *redis.Pool
}

func NewJsonService(prefix string, redisPool *redis.Pool) *JsonService {
	return &JsonService{prefix, redisPool}
}

func (s *JsonService) GetAll() (string, error) {
	conn := s.redisPool.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	keys, err := redis.Strings(conn.Do("KEYS", s.prefix+"*"))
	if err != nil {
		return "", err
	}

	result := strings.Builder{}
	result.WriteString("[")
	for i, key := range keys {
		if i > 0 {
			result.WriteString(",")
		}
		data, err := redis.String(conn.Do("GET", key))
		if err != nil {
			return "", err
		}

		result.WriteString(data)
	}
	result.WriteString("]")
	return result.String(), nil
}

func (s *JsonService) GetById(id string) (string, error) {
	conn := s.redisPool.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	key := s.prefix + id
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
