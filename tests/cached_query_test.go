package tests

func (suite *IntegrationTestSuite) TestCachedQueryFilterOneElementByString() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.2, Ok: false}
	suite.PutToRedisAsJson("cached-query.1", originalQuery1)
	suite.PutToRedisAsJson("cached-query.2", originalQuery2)
	suite.WaitForCacheDuration()

	// when
	var result []Query
	suite.HttpGetJson("/cached-query?Name=TestQuery", &result)

	// then
	suite.Equal(1, len(result))
	suite.Equal(originalQuery1, result[0])
}

func (suite *IntegrationTestSuite) TestCachedQueryFilterOneElementByInt() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.2, Ok: false}
	suite.PutToRedisAsJson("cached-query.1", originalQuery1)
	suite.PutToRedisAsJson("cached-query.2", originalQuery2)
	suite.WaitForCacheDuration()

	// when
	var result []Query
	suite.HttpGetJson("/cached-query?Number=1", &result)

	// then
	suite.Equal(1, len(result))
	suite.Equal(originalQuery1, result[0])
}

func (suite *IntegrationTestSuite) TestCachedQueryFilterOneElementByFloat() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.2, Ok: false}
	suite.PutToRedisAsJson("cached-query.1", originalQuery1)
	suite.PutToRedisAsJson("cached-query.2", originalQuery2)
	suite.WaitForCacheDuration()

	// when
	var result []Query
	suite.HttpGetJson("/cached-query?FloatNumber=1.1", &result)

	// then
	suite.Equal(1, len(result))
	suite.Equal(originalQuery1, result[0])
}

func (suite *IntegrationTestSuite) TestCachedQueryFilterOneElementByBool() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.2, Ok: false}
	suite.PutToRedisAsJson("cached-query.1", originalQuery1)
	suite.PutToRedisAsJson("cached-query.2", originalQuery2)
	suite.WaitForCacheDuration()

	// when
	var result []Query
	suite.HttpGetJson("/cached-query?Ok=true", &result)

	// then
	suite.Equal(1, len(result))
	suite.Equal(originalQuery1, result[0])
}

func (suite *IntegrationTestSuite) TestCachedQueryFilterComplexFilters() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.1, Ok: false}
	originalQuery3 := Query{ID: "3", Name: "TestQuery2", Number: 2, FloatNumber: 2.1, Ok: false}
	originalQuery4 := Query{ID: "4", Name: "TestQuery4", Number: 2, FloatNumber: 2.2, Ok: false}
	suite.PutToRedisAsJson("cached-query.1", originalQuery1)
	suite.PutToRedisAsJson("cached-query.2", originalQuery2)
	suite.PutToRedisAsJson("cached-query.3", originalQuery3)
	suite.PutToRedisAsJson("cached-query.4", originalQuery4)
	suite.WaitForCacheDuration()

	// when
	var result []Query
	suite.HttpGetJson("/cached-query?Ok=false&FloatNumber=2.1", &result)

	// then
	suite.Equal(2, len(result))
	suite.Contains(result, originalQuery2)
	suite.Contains(result, originalQuery3)
}

func (suite *IntegrationTestSuite) TestCachedQuerySubFilterComplexFilters() {
	// given
	originalQuery1 := Query{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalQuery2 := Query{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.1, Ok: false}
	originalQuery3 := Query{ID: "3", Name: "TestQuery2", Number: 2, FloatNumber: 2.1, Ok: false}
	originalQuery4 := Query{ID: "4", Name: "TestQuery4", Number: 2, FloatNumber: 2.2, Ok: false}

	originalComplexQuery1 := ComplexQuery{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true, SubQuery: originalQuery1}
	originalComplexQuery2 := ComplexQuery{ID: "2", Name: "TestQuery2", Number: 2, FloatNumber: 2.1, Ok: false, SubQuery: originalQuery2}
	originalComplexQuery3 := ComplexQuery{ID: "3", Name: "TestQuery2", Number: 2, FloatNumber: 2.1, Ok: false, SubQuery: originalQuery3}
	originalComplexQuery4 := ComplexQuery{ID: "4", Name: "TestQuery4", Number: 2, FloatNumber: 2.2, Ok: false, SubQuery: originalQuery4}

	suite.PutToRedisAsJson("cached-complex-query.1", originalComplexQuery1)
	suite.PutToRedisAsJson("cached-complex-query.2", originalComplexQuery2)
	suite.PutToRedisAsJson("cached-complex-query.3", originalComplexQuery3)
	suite.PutToRedisAsJson("cached-complex-query.4", originalComplexQuery4)
	suite.WaitForCacheDuration()

	// when
	var result []ComplexQuery
	suite.HttpGetJson("/cached-complex-query?SubQuery.Ok=false&SubQuery.FloatNumber=2.1", &result)

	// then
	suite.Equal(2, len(result))
	suite.Contains(result, originalComplexQuery2)
	suite.Contains(result, originalComplexQuery3)
}

func (suite *IntegrationTestSuite) TestCachedQuerySubFilterComplexFilter() {
	// given
	originalQuery := Query{ID: "1", Name: "ComplexTestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalComplexQuery := ComplexQuery{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true, SubQuery: originalQuery}

	suite.PutToRedisAsJson("cached-complex-query.1", originalComplexQuery)
	suite.WaitForCacheDuration()

	// when
	var result []ComplexQuery
	suite.HttpGetJson("/cached-complex-query?Ok=true&SubQuery.Name=ComplexTestQuery", &result)

	// then
	suite.Equal(1, len(result))
	suite.Contains(result, originalComplexQuery)
}

func (suite *IntegrationTestSuite) TestCachedQuerySubFilterComplexFilterNotFound() {
	// given
	originalQuery := Query{ID: "1", Name: "ComplexTestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalComplexQuery := ComplexQuery{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: true, SubQuery: originalQuery}

	suite.PutToRedisAsJson("cached-complex-query.1", originalComplexQuery)
	suite.WaitForCacheDuration()

	// when
	var result []ComplexQuery
	suite.HttpGetJson("/cached-complex-query?Ok=false&SubQuery.Name=ComplexTestQuery", &result)

	// then
	suite.Equal(0, len(result))
}

func (suite *IntegrationTestSuite) TestCachedQuerySuperComplexSubFilter() {
	// given
	originalQuery := Query{ID: "1", Name: "ComplexTestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalComplexQuery := ComplexQuery{ID: "1", Name: "TestQuery", Number: 42, FloatNumber: 1.1, Ok: false, SubQuery: originalQuery}
	originalSuperComplexQuery := SuperComplexQuery{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: false, SuperComplexQuery: originalComplexQuery}
	originalSuperSuperComplexQuery := SuperSuperComplexQuery{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: false, SuperSuperComplexQuery: originalSuperComplexQuery}

	suite.PutToRedisAsJson("cached-super-complex-query.1", originalSuperSuperComplexQuery)
	suite.WaitForCacheDuration()

	// when
	var result []SuperSuperComplexQuery
	suite.HttpGetJson("/cached-super-complex-query?SuperSuperComplexQuery.SuperComplexQuery.SubQuery.Ok=true&SuperSuperComplexQuery.SuperComplexQuery.Number=42", &result)

	// then
	suite.Equal(1, len(result))
	suite.Contains(result, originalSuperSuperComplexQuery)
}

func (suite *IntegrationTestSuite) TestCachedQuerySuperComplexSubFilterNotFound() {
	// given
	originalQuery := Query{ID: "1", Name: "ComplexTestQuery", Number: 1, FloatNumber: 1.1, Ok: true}
	originalComplexQuery := ComplexQuery{ID: "1", Name: "TestQuery", Number: 42, FloatNumber: 1.1, Ok: false, SubQuery: originalQuery}
	originalSuperComplexQuery := SuperComplexQuery{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: false, SuperComplexQuery: originalComplexQuery}
	originalSuperSuperComplexQuery := SuperSuperComplexQuery{ID: "1", Name: "TestQuery", Number: 1, FloatNumber: 1.1, Ok: false, SuperSuperComplexQuery: originalSuperComplexQuery}

	suite.PutToRedisAsJson("cached-super-complex-query.1", originalSuperSuperComplexQuery)
	suite.WaitForCacheDuration()

	// when
	var result []SuperSuperComplexQuery
	suite.HttpGetJson("/cached-super-complex-query?SuperSuperComplexQuery.SuperComplexQuery.SubQuery.Ok=false&SuperSuperComplexQuery.SuperComplexQuery.FloatNumber=1.0", &result)

	// then
	suite.Equal(0, len(result))
}
