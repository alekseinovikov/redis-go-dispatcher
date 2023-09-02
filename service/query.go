package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type QueryService struct {
	logger Logger
}

type Logger interface {
	Warn(i ...interface{})
	Error(i ...interface{})
}

func NewQueryService(logger Logger) *QueryService {
	return &QueryService{logger: logger}
}

func (s *QueryService) ApplyQuery(queryParams map[string][]string, data []string) ([]string, error) {
	if queryParams == nil || len(queryParams) == 0 {
		return data, nil
	}

	filters := s.buildFilters(queryParams)
	for _, filter := range filters {
		data = filter.Apply(data)
	}

	return data, nil
}

func (s *QueryService) buildFilters(queryParams map[string][]string) []filter {
	filters := make([]filter, 0, len(queryParams))
	for key, values := range queryParams {
		filters = append(filters, filter{
			logger:    s.logger,
			fieldPath: buildPath(key),
			values:    values,
		})
	}

	return filters
}

type filter struct {
	logger    Logger
	fieldPath path
	values    []string
}

func (f *filter) Apply(data []string) []string {
	if len(data) == 0 {
		return data
	}

	result := make([]string, 0, len(data))
	for _, jsonString := range data {
		jsonMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(jsonString), &jsonMap)
		if err != nil {
			f.logger.Error(err)
			continue
		}

		if val, ok := f.fieldPath.findValue(jsonMap); ok {
			valueString := convertToString(val)

			if Contains(f.values, valueString) {
				result = append(result, jsonString)
			}
		}
	}
	return result
}

func convertToString(val interface{}) string {
	valueString := ""
	if v, ok := val.(string); ok {
		valueString = v
	}
	if v, ok := val.(float64); ok {
		valueString = strconv.FormatFloat(v, 'f', -1, 64)
	}
	if v, ok := val.(bool); ok {
		valueString = fmt.Sprintf("%t", v)
	}
	if v, ok := val.(int); ok {
		valueString = fmt.Sprintf("%d", v)
	}
	return valueString
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

type path struct {
	field string
	next  *path
}

func (p *path) findValue(values map[string]interface{}) (string, bool) {
	stepValue, found := values[p.field]
	if !found {
		return "", false
	}

	if p.next == nil {
		return convertToString(stepValue), true
	}

	if stepValues, ok := stepValue.(map[string]interface{}); ok {
		return p.next.findValue(stepValues)
	}

	return "", false
}

func buildPath(field string) path {
	subFields := strings.Split(field, ".")
	return *createPath(subFields[0], subFields[1:])
}

func createPath(field string, rest []string) *path {
	if len(rest) == 0 {
		return &path{field: field}
	}
	return &path{field: field, next: createPath(rest[0], rest[1:])}
}
