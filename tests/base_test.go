package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	rt "github.com/testcontainers/testcontainers-go/modules/redis"
	"net/http"
	. "redis-go-dispatcher/config"
	"redis-go-dispatcher/server"
	"strconv"
	"testing"
	"time"
)

type IntegrationTestSuite struct {
	suite.Suite
	URLPrefix      string
	RedisPool      *redis.Pool
	RedisContainer *rt.RedisContainer
}

var cacheDuration = 300 * time.Millisecond

func testPrefixes() []Prefix {
	return []Prefix{
		{
			URI:                  "/cars",
			RedisPrefix:          "cars.",
			CacheEnabled:         true,
			CacheRefreshDuration: cacheDuration,
			CacheTtl:             cacheDuration,
		}, {
			URI:                  "/people",
			RedisPrefix:          "people.",
			CacheEnabled:         true,
			CacheRefreshDuration: cacheDuration,
			CacheTtl:             cacheDuration,
		},
	}
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	redisContainer := suite.startRedisContainer(ctx)
	connectionString := suite.createRedisConnPool(redisContainer, ctx)
	suite.startWebServer(connectionString)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	_ = suite.RedisPool.Close()
	_ = suite.RedisContainer.Terminate(context.Background())
}

func (suite *IntegrationTestSuite) TearDownTest() {
	conn := suite.RedisPool.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	_, _ = conn.Do("FLUSHALL")
	suite.WaitForCacheDuration()
}

func (suite *IntegrationTestSuite) startWebServer(connectionString string) {
	port := getFreePort()
	go server.StartServer(Config{
		Prefixes: testPrefixes(),
		Redis: RedisConfig{
			URL:           connectionString,
			PoolMaxIdle:   5,
			PoolMaxActive: 10,
		},
		ServerPort: port,
	})

	// Waiting for the server to start
	time.Sleep(1 * time.Second)
	suite.URLPrefix = fmt.Sprintf("http://localhost:%s", port)
}

func (suite *IntegrationTestSuite) createRedisConnPool(redisContainer *rt.RedisContainer, ctx context.Context) string {
	connectionString, err := redisContainer.ConnectionString(ctx)
	require.NoError(suite.T(), err)

	suite.RedisPool = &redis.Pool{
		MaxIdle:   5,
		MaxActive: 10,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(connectionString)
		},
	}
	return connectionString
}

func (suite *IntegrationTestSuite) startRedisContainer(ctx context.Context) *rt.RedisContainer {
	redisContainer, err := rt.RunContainer(ctx)
	require.NoError(suite.T(), err)
	suite.RedisContainer = redisContainer
	return redisContainer
}

func getFreePort() string {
	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(port)
}

func (suite *IntegrationTestSuite) WaitForCacheDuration() {
	time.Sleep(cacheDuration + 200*time.Millisecond)
}

func (suite *IntegrationTestSuite) PutToRedisAsJson(key string, obj interface{}) {
	conn := suite.RedisPool.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	value, _ := json.Marshal(obj)
	_, _ = conn.Do("SET", key, value)
	suite.WaitForCacheDuration()
}

func (suite *IntegrationTestSuite) DeleteFromRedis(key string) {
	conn := suite.RedisPool.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	_, _ = conn.Do("DEL", key)
	suite.WaitForCacheDuration()
}

func (suite *IntegrationTestSuite) HttpGetJson(uri string, target interface{}) {
	resp, err := http.Get(suite.URLPrefix + uri)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(target)
	_ = resp.Body.Close()
	assert.NoError(suite.T(), err)
}

func (suite *IntegrationTestSuite) HttpGet(uri string) *http.Response {
	resp, err := http.Get(suite.URLPrefix + uri)
	assert.NoError(suite.T(), err)
	return resp
}
