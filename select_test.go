package fselect

import "testing"

const errInvalidValue = "value is not a struct or pointer to struct"

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

type Pet struct {
	FirstName string `col:"first_name"`
	LastName  string `col:"last_name"`
	Age       int    `col:"age"`
}

func TestInvalidVAll(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrInvalidV {
			t.Fatal("Expected ErrInvalidV")
		}
	}()

	All(0)
}

func TestInvalidVAllExcept(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrInvalidV {
			t.Fatal("Expected ErrInvalidV")
		}
	}()

	AllExcept(0)
}

func TestInvalidVOnly(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrInvalidV {
			t.Fatal("Expected ErrInvalidV")
		}
	}()

	Only(0)
}

func TestFieldsNotFoundAllExcept(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrSomeFieldsNotFound {
			t.Fatal("Expected ErrSomeFieldsNotFound")
		}
	}()

	AllExcept(newPerson(), "this is an invalid field")
}

func TestFieldsNotFoundOnly(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrSomeFieldsNotFound {
			t.Fatal("Expected ErrSomeFieldsNotFound")
		}
	}()

	Only(newPerson(), "this is an invalid field")
}

func TestFields(t *testing.T) {
	p := newPerson()
	s := All(p)
	if s.Fields()[2] != "Age" {
		t.Fatal(`assert: s.Fields()[2] != "Age"`)
	}
}

func TestFieldNames(t *testing.T) {
	p := newPet()
	s := AllExcept(p, "first_name")
	if s.FieldString() != "last_name, age" {
		t.Fatal(`assert: s.FieldString() != "last_name, age"`)
	}
}

func TestBindVars(t *testing.T) {
	p := newPerson()
	s := Only(p, "Age")
	if s.BindVars() != "?" {
		t.Fatal(`assert: s.BindVars() != "?"`)
	}
}

func TestArgs(t *testing.T) {
	p := newPet()
	s := Only(p, "last_name")
	if s.Args()[0] != p.LastName {
		t.Fatal(`assert: s.Args()[0] != p.LastName`)
	}
}

func TestArgsAnd(t *testing.T) {
	additional := "additional_value"
	p := newPerson()
	s := All(p)
	args := s.ArgsAnd(additional)
	if args[len(args)-1] != additional {
		t.Fatal(`assert: args[len(args)-1] != additional`)
	}
}

func TestPrepare(t *testing.T) {
	const expect = "INSERT INTO pets (first_name, last_name, age) VALUES (?, ?, ?)"
	query := All(newPet()).Prepare("INSERT INTO pets (%fields%) VALUES (%vars%)")

	if query != expect {
		t.Fatal(`assert: query != expect`)
	}
}

func TestPreparePartial(t *testing.T) {
	const expect = "SELECT (first_name, last_name, age) FROM pets"
	query := All(newPet()).Prepare("SELECT (%fields%) FROM pets")

	if query != expect {
		t.Fatal(`assert: query != expect`)
	}
}
