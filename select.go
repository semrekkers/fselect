// Package fselect provides a simple struct field selector for preparing SQL queries.
package fselect

import (
	"errors"
	"reflect"
	"strings"
)

const (
	// StructTagKey is the struct field tag key.
	StructTagKey = "col"

	// FieldSeperator is the seperator between fields in a string.
	FieldSeperator = ", "

	// BindVar is the bind var to use. Only this MySQL bind var is supported right now.
	BindVar = "?"
)

const (
	fieldsVerb = "%fields%"
	varsVerb   = "%vars%"
	ignoreCase = true
)

var (
	// ErrInvalidV means that the argument v is not a struct or a pointer to a struct.
	ErrInvalidV = errors.New("v is not a struct or pointer to struct")

	// ErrSomeFieldsNotFound means that some fields you want to select where not found. Example:
	// fselect.AllExcept(&MyStruct{}, "a field that does not exists in struct MyStruct")
	//                                ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
	ErrSomeFieldsNotFound = errors.New("some fields wheren't found, check your field selection")
)

// Selection contains the selected fields.
type Selection struct {
	fieldNames  []string
	fieldValues []interface{}
}

// All selects all fields of struct v.
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

// AllExcept selects all fields of struct v except the fields specified in variadic argument fields. Differing cases are ignored.
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

		if sliceContains(fieldName, ignoreCase, fields) {
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

// Only selects only the fields of struct v specified in variadic argument fields. Differing cases are ignored.
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
		if !sliceContains(fieldName, ignoreCase, fields) {
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

// Fields returns all names of the selected fields.
func (s *Selection) Fields() []string {
	return s.fieldNames
}

// Args returns the values of all selected fields in order.
func (s *Selection) Args() []interface{} {
	return s.fieldValues
}

// FieldString returns all names of the selected fields seperated by const FieldSeperator.
func (s *Selection) FieldString() string {
	return strings.Join(s.fieldNames, FieldSeperator)
}

// BindVars returns const BindVar repeated n times where n is the amount of fields.
func (s *Selection) BindVars() string {
	return repeatString(BindVar, FieldSeperator, len(s.fieldNames))
}

// Prepare prepares qeury q. Two verbs will be replaced multiple times:
//    %fields% will be replaced with FieldString()
//    %vars% will be replaced with BindVars()
func (s *Selection) Prepare(q string) string {
	prepared := strings.Replace(q, fieldsVerb, s.FieldString(), -1)
	return strings.Replace(prepared, varsVerb, s.BindVars(), -1)
}
