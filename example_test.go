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
