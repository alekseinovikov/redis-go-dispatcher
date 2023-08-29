package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"time"
)

type Car struct {
	ID    string
	Model string
	Year  int
}

func (suite *IntegrationTestSuite) TestGetById() {
	// given
	original := Car{ID: "1", Model: "Toyota", Year: 2022}
	suite.PutToRedisAsJson("cars.1", original)

	// when
	var result Car
	suite.HttpGetJson("/cars/1", &result)

	// then
	assert.Equal(suite.T(), original, result)
}

func (suite *IntegrationTestSuite) TestGetAllCars() {
	conn := suite.RedisPool.Get()
	defer conn.Close()

	car1 := Car{ID: "1", Model: "Toyota", Year: 2022}
	car2 := Car{ID: "2", Model: "Honda", Year: 2023}

	data1, _ := json.Marshal(car1)
	data2, _ := json.Marshal(car2)

	conn.Do("SET", "cars.1", data1)
	conn.Do("SET", "cars.2", data2)

	time.Sleep(1 * time.Second)

	resp, err := http.Get(suite.URLPrefix + "/cars")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var cars []Car
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&cars)
	resp.Body.Close()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), cars, 2)

	resp, err = http.Get(suite.URLPrefix + "/cars/1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var car Car
	decoder = json.NewDecoder(resp.Body)
	err = decoder.Decode(&car)
	resp.Body.Close()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "1", car.ID)
	assert.Equal(suite.T(), "Toyota", car.Model)
	assert.Equal(suite.T(), 2022, car.Year)
}
