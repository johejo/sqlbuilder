package sqlbuilder_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/johejo/sqlbuilder"
)

func TestBuilder_Build(t *testing.T) {
	want := "SELECT * FROM test_user WHERE id = ? AND name = ? ORDER BY name LIMIT ?, ?"
	wantArgs := []interface{}{1, "John", 0, 10}

	var b sqlbuilder.Builder
	b.Append("SELECT * FROM test_user WHERE id = ? ", 1) // with trailing space
	b.Append("AND name = ?", "John")                     // without trailing space
	b.Append(" ORDER BY name LIMIT ?, ?", 0, 10)         // with space at the beginning

	got, gotArgs := b.Build()
	assertQuery(t, want, got)
	assertArgs(t, wantArgs, gotArgs)
}

func TestBuilder_Reset(t *testing.T) {
	var b sqlbuilder.Builder
	b.Append("SELECT * FROM test_user WHERE id = ?", 1)
	b.Reset()

	got, gotArgs := b.Build()
	if got != "" {
		t.Errorf("failed to reset builder: got=%s", got)
	}
	if len(gotArgs) != 0 {
		t.Errorf("failed to reset args: got=%v", gotArgs)
	}
}

func TestNewBulkBuilder(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "ok",
			query:   "INSERT INTO tbl_name (a,b,c) VALUES (?)",
			wantErr: false,
		},
		{
			name:    "error: no bulk insert marker",
			query:   "INSERT IGNORE INTO tbl_name (a,b,c) VALUES (?,?,?)",
			wantErr: true,
		},
		{
			name:    "error: not end with bulk insert marker",
			query:   "INSERT INTO tbl_name (a,b,c) VALUES (?),",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := sqlbuilder.NewBulkBuilder(tt.query); (err != nil) != tt.wantErr {
				t.Fatal(err)
			}
		})
	}
}

func TestBulkBuilder_Build(t *testing.T) {
	newTestBulkBuilder := func(t *testing.T) *sqlbuilder.BulkBuilder {
		t.Helper()
		const q = "INSERT INTO tbl_name (a,b,c) VALUES (?)"
		b, err := sqlbuilder.NewBulkBuilder(q)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}

	t.Run("ok", func(t *testing.T) {
		tests := []struct {
			name      string
			args      [][]interface{}
			wantQuery string
			wantArgs  []interface{}
		}{
			{
				name:      "insert one",
				args:      [][]interface{}{{"a", "b", "c"}},
				wantQuery: "INSERT INTO tbl_name (a,b,c) VALUES (?,?,?)",
				wantArgs:  []interface{}{"a", "b", "c"},
			},
			{
				name:      "bulk insert",
				args:      [][]interface{}{{"a1", "b1", "c1"}, {"a2", "b2", "c2"}},
				wantQuery: "INSERT INTO tbl_name (a,b,c) VALUES (?,?,?),(?,?,?)",
				wantArgs:  []interface{}{"a1", "b1", "c1", "a2", "b2", "c2"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				b := newTestBulkBuilder(t)
				for _, args := range tt.args {
					if err := b.Bind(args...); err != nil {
						t.Fatal(err)
					}
				}
				q, args := b.Build()
				assertQuery(t, tt.wantQuery, q)
				assertArgs(t, tt.wantArgs, args)
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		tests := []struct {
			name string
			args [][]interface{}
		}{
			{
				name: "invalid number of args",
				args: [][]interface{}{{"a", "b", "c", "d"}},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				b := newTestBulkBuilder(t)
				for _, args := range tt.args {
					if err := b.Bind(args...); err == nil {
						t.Fatal(err)
					}
				}
			})
		}
	})

	t.Run("empty", func(t *testing.T) {
		b := newTestBulkBuilder(t)
		if err := b.Bind(); err != nil {
			t.Fatal(err)
		}
	})
}

func assertQuery(t *testing.T, want, got string) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("failed to build query\n" + diff)
	}
}

func assertArgs(t *testing.T, want, got []interface{}) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("failed to build args: want=%v, got=%v", want, got)
	}
}
