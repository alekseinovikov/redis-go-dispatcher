package tests

import (
	"github.com/stretchr/testify/assert"
)

type Car struct {
	ID    string
	Model string
	Year  int
}

func (suite *IntegrationTestSuite) TestGetByIdFound() {
	// given
	original := Car{ID: "1", Model: "Toyota", Year: 2022}
	suite.PutToRedisAsJson("cars.1", original)

	// when
	var result Car
	suite.HttpGetJson("/cars/1", &result)

	// then
	assert.Equal(suite.T(), original, result)
}

func (suite *IntegrationTestSuite) TestGetByIdNotFound() {
	// when
	response := suite.HttpGet("/cars/1")

	// then
	assert.Equal(suite.T(), 404, response.StatusCode)
}

func (suite *IntegrationTestSuite) TestGetAllFound() {
	// given
	car1 := Car{ID: "1", Model: "Toyota", Year: 2022}
	car2 := Car{ID: "2", Model: "Honda", Year: 2023}
	suite.PutToRedisAsJson("cars.1", car1)
	suite.PutToRedisAsJson("cars.2", car2)

	// when
	var result []Car
	suite.HttpGetJson("/cars", &result)

	// then
	assert.Equal(suite.T(), 2, len(result))
	assert.Contains(suite.T(), result, car1)
	assert.Contains(suite.T(), result, car2)
}

func (suite *IntegrationTestSuite) TestGetAllNoElementsReturnsEmptyList() {
	// when
	var result []Car
	suite.HttpGetJson("/cars", &result)

	// then
	assert.Equal(suite.T(), 0, len(result))
}
