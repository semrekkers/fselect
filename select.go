package fselect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	StructTagKey   = "col"
	FieldSeperator = ","
	BindVar        = "?"
)

var (
	ErrInvalidV           = errors.New("v is not a struct or pointer to struct")
	ErrSomeFieldsNotFound = errors.New("some fields wheren't found, check your field selection")
)

type Selection struct {
	fieldNames  []string
	fieldValues []interface{}
}

func All(v interface{}) *Selection {
	value := reflect.Indirect(reflect.ValueOf(v))
	if value.Kind() != reflect.Struct {
		panic(ErrInvalidV)
	}
	structType := value.Type()

	var s Selection
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldType := structType.Field(i)
		s.fieldNames = append(s.fieldNames, getFieldName(&fieldType))
		s.fieldValues = append(s.fieldValues, fieldValue.Interface())
	}

	return &s
}

func AllExcept(v interface{}, fields ...string) *Selection {
	value := reflect.Indirect(reflect.ValueOf(v))
	if value.Kind() != reflect.Struct {
		panic(ErrInvalidV)
	}
	structType := value.Type()

	var s Selection
	for i := 0; i < value.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldName := getFieldName(&fieldType)

		if sliceContains(fieldName, fields) {
			// Skip current field, goto next
			continue
		}

		fieldValue := value.Field(i)
		s.fieldNames = append(s.fieldNames, fieldName)
		s.fieldValues = append(s.fieldValues, fieldValue.Interface())
	}

	if len(s.fieldNames) != value.NumField()-len(fields) {
		panic(ErrSomeFieldsNotFound)
	}

	return &s
}

func Only(v interface{}, fields ...string) *Selection {
	value := reflect.Indirect(reflect.ValueOf(v))
	if value.Kind() != reflect.Struct {
		panic(ErrInvalidV)
	}
	structType := value.Type()

	var s Selection
	for i := 0; i < value.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldName := getFieldName(&fieldType)

		// NOTE: the ! before sliceContains
		if !sliceContains(fieldName, fields) {
			// Skip current field, goto next
			continue
		}

		fieldValue := value.Field(i)
		s.fieldNames = append(s.fieldNames, fieldName)
		s.fieldValues = append(s.fieldValues, fieldValue.Interface())
	}

	if len(s.fieldNames) != len(fields) {
		panic(ErrSomeFieldsNotFound)
	}

	return &s
}

func (s *Selection) Fields() []string {
	return s.fieldNames
}

func (s *Selection) Args() []interface{} {
	return s.fieldValues
}

func (s *Selection) FieldString() string {
	return strings.Join(s.fieldNames, FieldSeperator)
}

func (s *Selection) BindVars() string {
	return repeatString(BindVar, FieldSeperator, len(s.fieldNames))
}

func (s *Selection) Preparef(query string) string {
	return fmt.Sprintf(query, s.FieldString(), s.BindVars())
}
