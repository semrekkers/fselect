package qbuilder

import "testing"

func TestBindVars(t *testing.T) {
	const format = "$bindvars"
	const expect = "?, ?, ?, ?, ?"

	out := Select(testObj{}).Fmt(format)
	if out != expect {
		t.Errorf("failed, got: %#v, want: %#v", out, expect)
	}
}

func TestNames(t *testing.T) {
	const format = "$names"
	const expect = "ID, User, Password, user_bio, user_website"

	out := Select(testObj{}).Fmt(format)
	if out != expect {
		t.Errorf("failed, got: %#v, want: %#v", out, expect)
	}
}

func TestUpdates(t *testing.T) {
	const format = "$updates"
	const expect = "ID = ?, User = ?, Password = ?, user_bio = ?, user_website = ?"

	out := Select(testObj{}).Fmt(format)
	if out != expect {
		t.Errorf("failed, got: %#v, want: %#v", out, expect)
	}
}

func TestTable(t *testing.T) {
	const format = "$table"
	const expect = "ID INT NOT NULL, User VARCHAR(255) NOT NULL, Password VARBINARY(255) NOT NULL, user_bio VARCHAR(255), user_website "

	out := Select(testObj{}).Fmt(format)
	if out != expect {
		t.Errorf("failed, got: %#v, want: %#v", out, expect)
	}
}
