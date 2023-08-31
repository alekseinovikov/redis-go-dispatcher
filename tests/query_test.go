package tests

type Query struct {
	ID          string
	Name        string
	Number      int
	FloatNumber float64
	Ok          bool
}

func (suite *IntegrationTestSuite) TestQueryFilterOneElementByString() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.2, Ok: false}
	suite.PutToRedisAsJson("query.1", originalQuery1)
	suite.PutToRedisAsJson("query.2", originalQuery2)

	// when
	var result []Query
	suite.HttpGetJson("/query?Name=TestQuery", &result)

	// then
	suite.Equal(1, len(result))
	suite.Equal(originalQuery1, result[0])
}

func (suite *IntegrationTestSuite) TestQueryFilterOneElementByInt() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.2, Ok: false}
	suite.PutToRedisAsJson("query.1", originalQuery1)
	suite.PutToRedisAsJson("query.2", originalQuery2)

	// when
	var result []Query
	suite.HttpGetJson("/query?Number=1", &result)

	// then
	suite.Equal(1, len(result))
	suite.Equal(originalQuery1, result[0])
}

func (suite *IntegrationTestSuite) TestQueryFilterOneElementByFloat() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.2, Ok: false}
	suite.PutToRedisAsJson("query.1", originalQuery1)
	suite.PutToRedisAsJson("query.2", originalQuery2)

	// when
	var result []Query
	suite.HttpGetJson("/query?FloatNumber=1.1", &result)

	// then
	suite.Equal(1, len(result))
	suite.Equal(originalQuery1, result[0])
}

func (suite *IntegrationTestSuite) TestQueryFilterOneElementByBool() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.2, Ok: false}
	suite.PutToRedisAsJson("query.1", originalQuery1)
	suite.PutToRedisAsJson("query.2", originalQuery2)

	// when
	var result []Query
	suite.HttpGetJson("/query?Ok=true", &result)

	// then
	suite.Equal(1, len(result))
	suite.Equal(originalQuery1, result[0])
}

func (suite *IntegrationTestSuite) TestQueryFilterComplexFilters() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.1, Ok: false}
	originalQuery3 := Query{ID: "3", Name: "TestQuery2", Number: 2, FloatNumber: 2.1, Ok: false}
	originalQuery4 := Query{ID: "4", Name: "TestQuery4", Number: 2, FloatNumber: 2.2, Ok: false}
	suite.PutToRedisAsJson("query.1", originalQuery1)
	suite.PutToRedisAsJson("query.2", originalQuery2)
	suite.PutToRedisAsJson("query.3", originalQuery3)
	suite.PutToRedisAsJson("query.4", originalQuery4)

	// when
	var result []Query
	suite.HttpGetJson("/query?Ok=false&FloatNumber=2.1", &result)

	// then
	suite.Equal(2, len(result))
	suite.Contains(result, originalQuery2)
	suite.Contains(result, originalQuery3)
}
