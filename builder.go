// Package sqlbuilder provides minimal functionality for building raw queries.
package sqlbuilder

import (
	"errors"
	"fmt"
	"strings"
)

// Builder is minimal SQL builder
type Builder struct {
	sb   strings.Builder
	args []interface{}
}

// NewBuilder returns a new *Builder.
func NewBuilder(opts ...Option) *Builder {
	cfg := new(config)
	for _, opt := range append(defaults(), opts...) {
		opt(cfg)
	}
	var sb strings.Builder
	return &Builder{
		sb:   sb,
		args: make([]interface{}, 0, cfg.size),
	}
}

// Append appends query and args to builder.
func (b *Builder) Append(q string, args ...interface{}) {
	b.appendQuery(q)
	b.sb.WriteString(" ")
	b.AppendArgs(args...)
}

func (b *Builder) appendQuery(q string) {
	q = strings.TrimSpace(q)
	b.sb.WriteString(q)
}

// AppendNoSpace appends query and args to builder without space.
// Using this method can cause SQL syntax errors.
func (b *Builder) AppendNoSpace(q string, args ...interface{}) {
	b.appendQuery(q)
	b.AppendArgs(args...)
}

// AppendArgs appends args.
func (b *Builder) AppendArgs(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	b.args = append(b.args, args...)
}

// Build returns the built query and args.
func (b *Builder) Build() (string, []interface{}) {
	return b.Query(), b.Args()
}

// Query returns the built query.
func (b *Builder) Query() string {
	return strings.TrimSuffix(b.sb.String(), " ")
}

// Args returns the built args.
func (b *Builder) Args() []interface{} {
	return b.args
}

// Reset resets internal buffer.
// Call Reset when you want to reuse the Builder.
func (b *Builder) Reset() {
	b.sb.Reset()
	b.args = make([]interface{}, 0)
}

// BulkBuilder is a minimal SQL builder for bulk insert.
// It is not supposed to be initialized by any method other than NewBulkBuilder.
type BulkBuilder struct {
	placeholders string
	n            int // number of placeholder per one bind
	builder      *Builder
}

// NewBulkBuilder is a minimal SQL builder for bulk insert.
func NewBulkBuilder(q string, opts ...Option) (*BulkBuilder, error) {
	cfg := new(config)
	opts = append(defaults(), opts...)
	for _, opt := range opts {
		opt(cfg)
	}
	q = strings.TrimSpace(q)
	if !strings.HasSuffix(q, cfg.marker) || strings.Count(q, cfg.marker) != 1 {
		return nil, fmt.Errorf("query must contains and end with only one %s as a marker for bulk builder", cfg.marker)
	}
	n := strings.Count(q, ",") + 1
	builder := NewBuilder(opts...)
	builder.Append(strings.TrimSuffix(q, cfg.marker))
	return &BulkBuilder{
		builder:      builder,
		n:            n,
		placeholders: fmt.Sprintf("(%s),", strings.TrimSuffix(strings.Repeat(cfg.placeholder+",", n), ",")),
	}, nil
}

// Bind binds args into BulkBuilder.
func (b *BulkBuilder) Bind(args ...interface{}) error {
	if len(args) == 0 {
		return nil
	}
	if b.n != len(args) {
		return errors.New(`number of placeholder and args is different`)
	}
	b.builder.AppendNoSpace(b.placeholders, args...)
	return nil
}

// Build returns the built query and args.
func (b *BulkBuilder) Build() (string, []interface{}) {
	return strings.TrimSuffix(b.builder.Query(), ","), b.builder.Args()
}

type config struct {
	size        int
	marker      string
	placeholder string
}

// Option is a option for NewBuilder and NewBulkBuilder.
type Option func(cfg *config)

// WithSize returns a Option to change the size used when initializing the internal slice.
// The default is 0.
// If specified negative number size will set 0.
func WithSize(size int) Option {
	if size < 0 {
		size = 0
	}
	return func(cfg *config) {
		cfg.size = size
	}
}

// WithMarker returns a Option to set the marker for use with bulk insert.
// The default is "(?)".
// Ignored if using with NewBuilder.
func WithMarker(marker string) Option {
	if marker == "" {
		marker = "(?)"
	}
	return func(cfg *config) {
		cfg.marker = marker
	}
}

// WithPlaceHolder returns a Option to set the marker for use with bulk insert.
// The default is "?".
// Ignored if using with NewBuilder.
func WithPlaceHolder(s string) Option {
	if s == "" {
		s = "?"
	}
	return func(cfg *config) {
		cfg.placeholder = s
	}
}

func defaults() []Option {
	return []Option{
		WithMarker(""),
		WithSize(-1),
		WithPlaceHolder(""),
	}
}
