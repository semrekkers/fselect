package qbuilder

import (
	"strings"
	"sync"
)

const (
	fmtBindVars = "$bindvars"
	fmtNames    = "$names"
	fmtUpdates  = "$updates"
	fmtTable    = "$table"
)

type formatOptions struct {
	fieldSeparator, bindVar string
}

// Formatter can format strings after a Selection is made.
type Formatter struct {
	fieldNames []string
	fieldTypes []string
	options    formatOptions

	lock             sync.Mutex
	compiledBindVars string
	compiledNames    string
	compiledUpdates  string
	compiledTable    string
}

// FieldNames returns the field names.
func (f *Formatter) FieldNames() []string {
	return f.fieldNames
}

// FieldTypes returns the field types.
func (f *Formatter) FieldTypes() []string {
	return f.fieldTypes
}

// BindVars returns a formatted string containing the bind var placeholders.
func (f *Formatter) BindVars() string {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.compiledBindVars == "" {
		f.compiledBindVars = repeatString(f.options.bindVar, f.options.fieldSeparator, len(f.fieldNames))
	}
	return f.compiledBindVars
}

// Names returns a formatted string containing the names of all the fields.
func (f *Formatter) Names() string {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.compiledNames == "" {
		f.compiledNames = strings.Join(f.fieldNames, f.options.fieldSeparator)
	}
	return f.compiledNames
}

// Updates returns a formatted string containing the field updates.
func (f *Formatter) Updates() string {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.compiledUpdates == "" {
		f.compiledUpdates = joinStringsWithSuffix(f.fieldNames, " = "+f.options.bindVar, f.options.fieldSeparator)
	}
	return f.compiledUpdates
}

// Table returns a formatted string containing the fields in table format.
func (f *Formatter) Table() string {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.compiledTable == "" {
		f.compiledTable = joinTwoSlices(f.fieldNames, f.fieldTypes, " ", f.options.fieldSeparator)
	}
	return f.compiledTable
}

// Fmt builds a string using format. The following strings are replaced:
//	$bindvars with the output of BindVars()
//	$names with the output of Names()
//	$updates with the output of Updates()
//	$table with the output of Table()
func (f *Formatter) Fmt(format string) string {
	if strings.Contains(format, fmtBindVars) {
		format = strings.Replace(format, fmtBindVars, f.BindVars(), -1)
	}
	if strings.Contains(format, fmtNames) {
		format = strings.Replace(format, fmtNames, f.Names(), -1)
	}
	if strings.Contains(format, fmtUpdates) {
		format = strings.Replace(format, fmtUpdates, f.Updates(), -1)
	}
	if strings.Contains(format, fmtTable) {
		format = strings.Replace(format, fmtTable, f.Table(), -1)
	}
	return format
}
