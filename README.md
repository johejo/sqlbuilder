# sqlbuilder

[![ci](https://github.com/johejo/sqlbuilder/workflows/ci/badge.svg?branch=master)](https://github.com/johejo/sqlbuilder/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/johejo/sqlbuilder/branch/master/graph/badge.svg)](https://codecov.io/gh/johejo/sqlbuilder)
[![Go Report Card](https://goreportcard.com/badge/github.com/johejo/sqlbuilder)](https://goreportcard.com/report/github.com/johejo/sqlbuilder)

## Description

Package sqlbuilder provides minimal functionality for building raw queries.

## Example

```go
package sqlbuilder_test

import (
	"fmt"

	"github.com/johejo/sqlbuilder"
)

func ExampleBuilder() {
	var b sqlbuilder.Builder
	b.Append("SELECT * FROM test_user WHERE id = ? ", 1) // with trailing space
	b.Append("AND name = ?", "John")                     // without trailing space
	b.Append(" ORDER BY name LIMIT ?, ?", 0, 10)         // with space at the beginning

	q, args := b.Build()
	fmt.Println(q)
	fmt.Println(args)

	// Output:
	// SELECT * FROM test_user WHERE id = ? AND name = ? ORDER BY name LIMIT ?, ?
	// [1 John 0 10]
}

func ExampleBulkBuilder() {
	b, err := sqlbuilder.NewBulkBuilder("INSERT INTO tbl_name (a,b,c) VALUES (?)")
	if err != nil {
		panic(err)
	}
	if err := b.Bind("A", "B", "C"); err != nil {
		panic(err)
	}
	if err := b.Bind("AA", "BB", "CC"); err != nil {
		panic(err)
	}

	q, args := b.Build()
	fmt.Println(q)
	fmt.Println(args)

	// Output:
	// INSERT INTO tbl_name (a,b,c) VALUES (?,?,?),(?,?,?)
	// [A B C AA BB CC]
}
```

## License

MIT

## Author

Mitsuo Heijo (@johejo)
