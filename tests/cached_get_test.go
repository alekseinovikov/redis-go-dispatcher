package tests

import (
	"github.com/stretchr/testify/assert"
)

type CachedCar struct {
	ID    string
	Model string
	Year  int
}

type CachedPerson struct {
	ID   string
	Name string
	Age  int
}

type CachedLongWarmUp struct {
	ID string
}

func (suite *IntegrationTestSuite) TestCachedGetByIdFound() {
	// given
	original := CachedCar{ID: "1", Model: "Toyota", Year: 2022}
	suite.PutToRedisAsJson("cached-cars.1", original)
	suite.WaitForCacheDuration()

	// when
	var result CachedCar
	suite.HttpGetJson("/cached-cars/1", &result)

	// then
	assert.Equal(suite.T(), original, result)
}

func (suite *IntegrationTestSuite) TestCachedGetByIdCacheIsNotWarmedUpNotFound() {
	// given
	original := CachedLongWarmUp{ID: "1"}
	suite.PutToRedisAsJson("cached-long_warm_up.1", original)

	// when
	response := suite.HttpGet("/cached-long_warm_up/1")

	// then
	assert.Equal(suite.T(), 404, response.StatusCode)
}

func (suite *IntegrationTestSuite) TestCachedGetByIdFoundAndThenDeleted() {
	// given
	original := CachedCar{ID: "1", Model: "Toyota", Year: 2022}
	suite.PutToRedisAsJson("cached-cars.1", original)
	suite.WaitForCacheDuration()

	var result CachedCar
	suite.HttpGetJson("/cached-cars/1", &result)
	assert.Equal(suite.T(), original, result)

	// when
	suite.DeleteFromRedis("cached-cars.1")
	suite.WaitForCacheDuration()

	response := suite.HttpGet("/cached-cars/1")

	// then
	assert.Equal(suite.T(), 404, response.StatusCode)
}

func (suite *IntegrationTestSuite) TestCachedGetByIdMultiplePrefixesFound() {
	// given
	originalCar := CachedCar{ID: "1", Model: "Toyota", Year: 2022}
	suite.PutToRedisAsJson("cached-cars.1", originalCar)

	originalPerson := CachedPerson{ID: "1", Name: "John", Age: 30}
	suite.PutToRedisAsJson("cached-people.1", originalPerson)

	suite.WaitForCacheDuration()

	// when
	var carResult CachedCar
	suite.HttpGetJson("/cached-cars/1", &carResult)
	var personResult CachedPerson
	suite.HttpGetJson("/cached-people/1", &personResult)

	// then
	assert.Equal(suite.T(), originalCar, carResult)
	assert.Equal(suite.T(), originalPerson, personResult)
}

func (suite *IntegrationTestSuite) TestCachedGetByIdMultiplePrefixesFoundAndThenOneDeleted() {
	// given
	originalCar := CachedCar{ID: "1", Model: "Toyota", Year: 2022}
	suite.PutToRedisAsJson("cached-cars.1", originalCar)

	originalPerson := CachedPerson{ID: "1", Name: "John", Age: 30}
	suite.PutToRedisAsJson("cached-people.1", originalPerson)

	suite.WaitForCacheDuration()

	var carResult CachedCar
	suite.HttpGetJson("/cached-cars/1", &carResult)
	var personResult CachedPerson
	suite.HttpGetJson("/cached-people/1", &personResult)

	assert.Equal(suite.T(), originalCar, carResult)
	assert.Equal(suite.T(), originalPerson, personResult)

	// when
	suite.DeleteFromRedis("cached-people.1")
	suite.WaitForCacheDuration()

	suite.HttpGetJson("/cached-cars/1", &carResult)
	response := suite.HttpGet("/cached-people/1")

	// then
	assert.Equal(suite.T(), 404, response.StatusCode)
	assert.Equal(suite.T(), originalCar, carResult)
}

func (suite *IntegrationTestSuite) TestCachedGetByIdNotFound() {
	// when
	response := suite.HttpGet("/cached-cars/1")

	// then
	assert.Equal(suite.T(), 404, response.StatusCode)
}

func (suite *IntegrationTestSuite) TestCachedGetByIdMultiplePrefixesNotFound() {
	// when
	responseCar := suite.HttpGet("/cached-cars/1")
	responsePerson := suite.HttpGet("/cached-people/1")

	// then
	assert.Equal(suite.T(), 404, responseCar.StatusCode)
	assert.Equal(suite.T(), 404, responsePerson.StatusCode)
}

func (suite *IntegrationTestSuite) TestCachedGetAllFound() {
	// given
	car1 := CachedCar{ID: "1", Model: "Toyota", Year: 2022}
	car2 := CachedCar{ID: "2", Model: "Honda", Year: 2023}
	suite.PutToRedisAsJson("cached-cars.1", car1)
	suite.PutToRedisAsJson("cached-cars.2", car2)

	suite.WaitForCacheDuration()

	// when
	var result []CachedCar
	suite.HttpGetJson("/cached-cars", &result)

	// then
	assert.Equal(suite.T(), 2, len(result))
	assert.Contains(suite.T(), result, car1)
	assert.Contains(suite.T(), result, car2)
}

func (suite *IntegrationTestSuite) TestCachedGetAllCacheIsNotWarmedUpEmptyList() {
	// given
	longWarmUp1 := CachedLongWarmUp{ID: "1"}
	longWarmUp2 := CachedLongWarmUp{ID: "2"}
	suite.PutToRedisAsJson("cached-long_warm_up.1", longWarmUp1)
	suite.PutToRedisAsJson("cached-long_warm_up.2", longWarmUp2)

	// when
	var result []CachedLongWarmUp
	suite.HttpGetJson("/cached-long_warm_up", &result)

	// then
	assert.Equal(suite.T(), 0, len(result))
}

func (suite *IntegrationTestSuite) TestCachedGetAllMultiplePrefixesFound() {
	// given
	car1 := CachedCar{ID: "1", Model: "Toyota", Year: 2022}
	car2 := CachedCar{ID: "2", Model: "Honda", Year: 2023}
	suite.PutToRedisAsJson("cached-cars.1", car1)
	suite.PutToRedisAsJson("cached-cars.2", car2)

	person1 := CachedPerson{ID: "1", Name: "John", Age: 30}
	person2 := CachedPerson{ID: "2", Name: "Jane", Age: 25}
	suite.PutToRedisAsJson("cached-people.1", person1)
	suite.PutToRedisAsJson("cached-people.2", person2)

	suite.WaitForCacheDuration()

	// when
	var resultCars []CachedCar
	suite.HttpGetJson("/cached-cars", &resultCars)
	var resultPeople []CachedPerson
	suite.HttpGetJson("/cached-people", &resultPeople)

	// then
	assert.Equal(suite.T(), 2, len(resultCars))
	assert.Contains(suite.T(), resultCars, car1)
	assert.Contains(suite.T(), resultCars, car2)

	assert.Equal(suite.T(), 2, len(resultPeople))
	assert.Contains(suite.T(), resultPeople, person1)
	assert.Contains(suite.T(), resultPeople, person2)
}

func (suite *IntegrationTestSuite) TestCachedGetAllNoElementsReturnsEmptyList() {
	// when
	var result []CachedCar
	suite.HttpGetJson("/cached-cars", &result)

	// then
	assert.Equal(suite.T(), 0, len(result))
}

func (suite *IntegrationTestSuite) TestCachedGetAllMultiplePrefixesNoElementsReturnsEmptyList() {
	// when
	var resultCars []CachedCar
	suite.HttpGetJson("/cached-cars", &resultCars)

	var resultPeople []CachedPerson
	suite.HttpGetJson("/cached-people", &resultPeople)

	// then
	assert.Equal(suite.T(), 0, len(resultCars))
	assert.Equal(suite.T(), 0, len(resultPeople))
}
