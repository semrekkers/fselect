package fselect

import "fmt"

func ExampleAll() {
	person := struct {
		Nickname  string
		FirstName string
		LastName  string
		Age       int
	}{}

	s := All(&person)
	fmt.Print(s.Fields())
	// Output: [Nickname FirstName LastName Age]
}

func ExampleAllExcept() {
	// Note that tags are used here.
	person := struct {
		Nickname  string `col:"nickname"`
		FirstName string `col:"first_name"`
		LastName  string `col:"last_name"`
		Age       int    `col:"age"`
	}{}

	// Don't select the Age field.
	s := AllExcept(&person, "age")

	fmt.Print(s.Fields())
	// Output: [nickname first_name last_name]
}

func ExampleOnly() {
	// Note that tags are used here.
	person := struct {
		Nickname  string `col:"nickname"`
		FirstName string `col:"first_name"`
		LastName  string `col:"last_name"`
		Age       int    `col:"age"`
	}{}

	// Only select the FirstName of person.
	s := Only(&person, "first_name")

	fmt.Print(s.Fields())
	// Output: [first_name]
}

func ExampleSelection_BindVars() {
	// Note that tags are used here.
	person := struct {
		Nickname  string `col:"nickname"`
		FirstName string `col:"first_name"`
		LastName  string `col:"last_name"`
		Age       int    `col:"age"`
	}{}

	s := All(&person)

	// Four fields, four times the BindVar.
	fmt.Print(s.BindVars())
	// Output: ?, ?, ?, ?
}

func ExampleSelection_Preparef() {
	person := struct {
		Nickname  string `col:"nickname"`
		FirstName string `col:"first_name"`
		LastName  string `col:"last_name"`
		Age       int    `col:"age"`
	}{}

	// Select all fields of person.
	s := All(&person)

	selectQuery := s.Preparef("SELECT (%fields%) FROM persons")
	insertQuery := s.Preparef("INSERT INTO persons (%fields%) VALUES (%vars%)")

	fmt.Println("select query:", selectQuery)
	fmt.Println("insert query:", insertQuery)

	// Output:
	// select query: SELECT (nickname, first_name, last_name, age) FROM persons
	// insert query: INSERT INTO persons (nickname, first_name, last_name, age) VALUES (?, ?, ?, ?)
}
