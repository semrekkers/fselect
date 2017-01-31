fselect
=======
A simple struct field selector to prepare SQL queries with a long list of fields.

### Why?
I found [jmoiron/sqlx](https://github.com/jmoiron/sqlx) a very useful package, but the only thing I missed was the ability to serialize struct fields into a SQL query. Luckily I developed this package to fill that gap.

### How to use
A nice example:
```go
// A struct with some fields
user := struct {
  ID         int
  Nickname   string
  FirstName  string
  LastName   string
  Email      string
  Password   []byte
  Location   string
  Website    string
  PrivateKey string
}{}

query := fselect.All(&user).Prepare("INSERT INTO users (%fields%) VALUES (%vars%)")
```

And `query` will be:
```sql
INSERT INTO users (ID, Nickname, FirstName, LastName, Email, Password, Location, Website, PrivateKey) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
```

You can also filter fields with `fselect.AllExcept()` and `fselect.Only()`.

### Docs
The docs are hosted on GoDoc.org [![GoDoc](https://godoc.org/github.com/semrekkers/fselect?status.svg)](https://godoc.org/github.com/semrekkers/fselect)

### Contributions
Contributions are more than welcome! Feel free to send pull requests and/or issues!
