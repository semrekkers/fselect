package qbuilder

import (
	"reflect"
	"testing"
)

type testObj struct {
	ID       string `type:"INT NOT NULL"`
	User     string `type:"VARCHAR(255) NOT NULL"`
	Password []byte `type:"VARBINARY(255) NOT NULL"`
	Bio      string `db:"user_bio" type:"VARCHAR(255)"`
	Website  string `db:"user_website"`
}

func TestSelect(t *testing.T) {
	obj := testObj{}
	s := Select(&obj)
	if s.structType != reflect.TypeOf(obj) {
		t.Error("failed")
	}
}

func TestSelectDefaults(t *testing.T) {
	s := Select(testObj{})
	if s.tagKey != defaultTagKey {
		t.Error("tagKey")
	}
	if s.typeTagKey != defaultTypeTagKey {
		t.Error("typeTagKey")
	}
	if s.options.fieldSeparator != defaultFieldSeparator {
		t.Error("fieldSeparator")
	}
	if s.options.bindVar != defaultBindVar {
		t.Error("bindVar")
	}
}

func TestFalseSelect(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrNotStructKind {
			t.Fatal("Expected ErrNotStructKind")
		}
	}()
	Select(0)
}

func TestExclude(t *testing.T) {
	s := Select(testObj{}).Exclude("ID")
	if !sliceContains("ID", s.filterSet) {
		t.Error("Field ID not found in filterSet")
	}
	if s.filterInclude {
		t.Error("filterInclude is set")
	}
}

func TestExcludeFail(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrFilterSet {
			t.Fatal("Expected ErrFilterSet")
		}
	}()
	Select(testObj{}).Exclude("ID").Exclude("User")
}

func TestOnly(t *testing.T) {
	s := Select(testObj{}).Only("ID")
	if !sliceContains("ID", s.filterSet) {
		t.Error("Field ID not found in filterSet")
	}
	if len(s.filterSet) != 1 {
		t.Error("Too many fields in filterSet")
	}
	if !s.filterInclude {
		t.Error("filterInclude is not set")
	}
}

func TestOnlyFail(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrFilterSet {
			t.Fatal("Expected ErrFilterSet")
		}
	}()
	Select(testObj{}).Only("ID").Only("User")
}

func TestBuildFormatter(t *testing.T) {
	fieldNames := []string{"ID", "User", "Password", "user_bio", "user_website"}
	fieldTypes := []string{"INT NOT NULL", "VARCHAR(255) NOT NULL", "VARBINARY(255) NOT NULL", "VARCHAR(255)", ""}

	f := Select(testObj{}).Formatter()
	if !sliceEqual(fieldNames, f.fieldNames) {
		t.Errorf("failed, got: %#v, want: %#v", f.fieldNames, fieldNames)
	}
	if !sliceEqual(fieldTypes, f.fieldTypes) {
		t.Errorf("failed, got: %#v, want: %#v", f.fieldTypes, fieldTypes)
	}
}

func TestBuildFormatterFilterExclude(t *testing.T) {
	fieldNames := []string{"User", "Password", "user_bio", "user_website"}
	fieldTypes := []string{"VARCHAR(255) NOT NULL", "VARBINARY(255) NOT NULL", "VARCHAR(255)", ""}

	f := Select(testObj{}).Exclude("ID").Formatter()
	if !sliceEqual(fieldNames, f.fieldNames) {
		t.Errorf("failed, got: %#v, want: %#v", f.fieldNames, fieldNames)
	}
	if !sliceEqual(fieldTypes, f.fieldTypes) {
		t.Errorf("failed, got: %#v, want: %#v", f.fieldTypes, fieldTypes)
	}
}

func TestBuildFormatterFilterOnly(t *testing.T) {
	fieldNames := []string{"User", "user_bio"}
	fieldTypes := []string{"VARCHAR(255) NOT NULL", "VARCHAR(255)"}

	f := Select(testObj{}).Only("User", "user_bio").Formatter()
	if !sliceEqual(fieldNames, f.fieldNames) {
		t.Errorf("failed, got: %#v, want: %#v", f.fieldNames, fieldNames)
	}
	if !sliceEqual(fieldTypes, f.fieldTypes) {
		t.Errorf("failed, got: %#v, want: %#v", f.fieldTypes, fieldTypes)
	}
}

func TestBuildFormatterFail(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("Expected err")
		}
	}()

	Select(testObj{}).Exclude("This doesn't exists").Formatter()
	Select(testObj{}).Only("This doesn't exists").Formatter()
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
