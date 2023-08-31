package tests

import (
	"github.com/stretchr/testify/assert"
)

type Car struct {
	ID    string
	Model string
	Year  int
}

type Person struct {
	ID   string
	Name string
	Age  int
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

func (suite *IntegrationTestSuite) TestGetByIdFoundAndThenDeleted() {
	// given
	original := Car{ID: "1", Model: "Toyota", Year: 2022}
	suite.PutToRedisAsJson("cars.1", original)

	var result Car
	suite.HttpGetJson("/cars/1", &result)
	assert.Equal(suite.T(), original, result)

	// when
	suite.DeleteFromRedis("cars.1")

	response := suite.HttpGet("/cars/1")

	// then
	assert.Equal(suite.T(), 404, response.StatusCode)
}

func (suite *IntegrationTestSuite) TestGetByIdMultiplePrefixesFound() {
	// given
	originalCar := Car{ID: "1", Model: "Toyota", Year: 2022}
	suite.PutToRedisAsJson("cars.1", originalCar)

	originalPerson := Person{ID: "1", Name: "John", Age: 30}
	suite.PutToRedisAsJson("people.1", originalPerson)

	// when
	var carResult Car
	suite.HttpGetJson("/cars/1", &carResult)
	var personResult Person
	suite.HttpGetJson("/people/1", &personResult)

	// then
	assert.Equal(suite.T(), originalCar, carResult)
	assert.Equal(suite.T(), originalPerson, personResult)
}

func (suite *IntegrationTestSuite) TestGetByIdMultiplePrefixesFoundAndThenOneDeleted() {
	// given
	originalCar := Car{ID: "1", Model: "Toyota", Year: 2022}
	suite.PutToRedisAsJson("cars.1", originalCar)

	originalPerson := Person{ID: "1", Name: "John", Age: 30}
	suite.PutToRedisAsJson("people.1", originalPerson)

	var carResult Car
	suite.HttpGetJson("/cars/1", &carResult)
	var personResult Person
	suite.HttpGetJson("/people/1", &personResult)

	assert.Equal(suite.T(), originalCar, carResult)
	assert.Equal(suite.T(), originalPerson, personResult)

	// when
	suite.DeleteFromRedis("people.1")

	suite.HttpGetJson("/cars/1", &carResult)
	response := suite.HttpGet("/people/1")

	// then
	assert.Equal(suite.T(), 404, response.StatusCode)
	assert.Equal(suite.T(), originalCar, carResult)
}

func (suite *IntegrationTestSuite) TestGetByIdNotFound() {
	// when
	response := suite.HttpGet("/cars/1")

	// then
	assert.Equal(suite.T(), 404, response.StatusCode)
}

func (suite *IntegrationTestSuite) TestGetByIdMultiplePrefixesNotFound() {
	// when
	responseCar := suite.HttpGet("/cars/1")
	responsePerson := suite.HttpGet("/people/1")

	// then
	assert.Equal(suite.T(), 404, responseCar.StatusCode)
	assert.Equal(suite.T(), 404, responsePerson.StatusCode)
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

func (suite *IntegrationTestSuite) TestGetAllMultiplePrefixesFound() {
	// given
	car1 := Car{ID: "1", Model: "Toyota", Year: 2022}
	car2 := Car{ID: "2", Model: "Honda", Year: 2023}
	suite.PutToRedisAsJson("cars.1", car1)
	suite.PutToRedisAsJson("cars.2", car2)

	person1 := Person{ID: "1", Name: "John", Age: 30}
	person2 := Person{ID: "2", Name: "Jane", Age: 25}
	suite.PutToRedisAsJson("people.1", person1)
	suite.PutToRedisAsJson("people.2", person2)

	// when
	var resultCars []Car
	suite.HttpGetJson("/cars", &resultCars)
	var resultPeople []Person
	suite.HttpGetJson("/people", &resultPeople)

	// then
	assert.Equal(suite.T(), 2, len(resultCars))
	assert.Contains(suite.T(), resultCars, car1)
	assert.Contains(suite.T(), resultCars, car2)

	assert.Equal(suite.T(), 2, len(resultPeople))
	assert.Contains(suite.T(), resultPeople, person1)
	assert.Contains(suite.T(), resultPeople, person2)
}

func (suite *IntegrationTestSuite) TestGetAllNoElementsReturnsEmptyList() {
	// when
	var result []Car
	suite.HttpGetJson("/cars", &result)

	// then
	assert.Equal(suite.T(), 0, len(result))
}

func (suite *IntegrationTestSuite) TestGetAllMultiplePrefixesNoElementsReturnsEmptyList() {
	// when
	var resultCars []Car
	suite.HttpGetJson("/cars", &resultCars)

	var resultPeople []Person
	suite.HttpGetJson("/people", &resultPeople)

	// then
	assert.Equal(suite.T(), 0, len(resultCars))
	assert.Equal(suite.T(), 0, len(resultPeople))
}
