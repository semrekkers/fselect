// Package qbuilder implements a simple, fast and easy-to-use query builder for jmoiron/sqlx.
package qbuilder

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrNotStructKind means that the given argument to Select() is not a struct kind.
	ErrNotStructKind = errors.New("argument is not a struct kind")

	// ErrFilterSet means that filters were already set.
	ErrFilterSet = errors.New("filters already set")
)

const (
	defaultTagKey         = "db"
	defaultTypeTagKey     = "type"
	defaultFieldSeparator = ", "
	defaultBindVar        = "?"
)

// Selection contains the selected fields and options.
type Selection struct {
	structType         reflect.Type
	filterSet          []string
	filterInclude      bool
	tagKey, typeTagKey string
	options            formatOptions
}

// Select selecs all the fields of s and returns a Selection.
func Select(s interface{}) *Selection {
	value := reflect.Indirect(reflect.ValueOf(s))
	if value.Kind() != reflect.Struct {
		panic(ErrNotStructKind)
	}
	return &Selection{
		structType: value.Type(),
		tagKey:     defaultTagKey,
		typeTagKey: defaultTypeTagKey,

		options: formatOptions{
			fieldSeparator: defaultFieldSeparator,
			bindVar:        defaultBindVar,
		},
	}
}

// Exclude excludes fields from Selection.
func (s *Selection) Exclude(fields ...string) *Selection {
	if s.filterSet != nil {
		panic(ErrFilterSet)
	}
	s.filterSet = fields
	return s
}

// Only excludes all the fields in the current Selection except fields.
func (s *Selection) Only(fields ...string) *Selection {
	if s.filterSet != nil {
		panic(ErrFilterSet)
	}
	s.filterSet = fields
	s.filterInclude = true
	return s
}

// TagKey sets the struct field tag key to use for this Selection.
func (s *Selection) TagKey(v string) *Selection {
	s.tagKey = v
	return s
}

// TypeTagKey sets the struct field tag key for types to use in this Selection.
func (s *Selection) TypeTagKey(v string) *Selection {
	s.typeTagKey = v
	return s
}

// FieldSeparator sets the field separator for the formatter.
func (s *Selection) FieldSeparator(v string) *Selection {
	s.options.fieldSeparator = v
	return s
}

// BindVar sets the bind variable placeholder for the formatter.
func (s *Selection) BindVar(v string) *Selection {
	s.options.bindVar = v
	return s
}

// Formatter builds and returns the Formatter for this Selection.
func (s *Selection) Formatter() *Formatter {
	return s.buildFormatter()
}

// Fmt builds the formatter for this Selection and returns a formatted string.
func (s *Selection) Fmt(format string) string {
	return s.buildFormatter().Fmt(format)
}

func (s *Selection) buildFormatter() (f *Formatter) {
	numField := s.structType.NumField()
	expectNumField := numField
	if s.filterSet != nil && s.filterInclude {
		expectNumField = len(s.filterSet)
	} else if s.filterSet != nil {
		expectNumField -= len(s.filterSet)
	}
	f = &Formatter{
		fieldNames: make([]string, 0, expectNumField),
		fieldTypes: make([]string, 0, expectNumField),
		options:    s.options,
	}

	for i := 0; i < numField; i++ {
		structField := s.structType.Field(i)
		fieldName := structField.Name
		if tagValue, ok := structField.Tag.Lookup(s.tagKey); ok {
			fieldName = tagValue
		}
		if s.filterSet != nil {
			fieldInSet := sliceContains(fieldName, s.filterSet)
			if (fieldInSet && !s.filterInclude) || (!fieldInSet && s.filterInclude) {
				// field is filtered out
				continue
			}
		}
		var fieldType string
		if tagValue, ok := structField.Tag.Lookup(s.typeTagKey); ok {
			fieldType = tagValue
		}
		f.fieldNames = append(f.fieldNames, fieldName)
		f.fieldTypes = append(f.fieldTypes, fieldType)
	}
	if len(f.fieldNames) != expectNumField {
		panic(fmt.Errorf("expected %d fields, have %d, please check your filters", expectNumField, len(f.fieldNames)))
	}
	return
}
