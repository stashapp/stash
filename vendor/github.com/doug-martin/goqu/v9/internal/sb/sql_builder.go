package sb

import (
	"bytes"
)

// Builder that is composed of a bytes.Buffer. It is used internally and by adapters to build SQL statements
type (
	SQLBuilder interface {
		Error() error
		SetError(err error) SQLBuilder
		WriteArg(i ...interface{}) SQLBuilder
		Write(p []byte) SQLBuilder
		WriteStrings(ss ...string) SQLBuilder
		WriteRunes(r ...rune) SQLBuilder
		IsPrepared() bool
		CurrentArgPosition() int
		ToSQL() (sql string, args []interface{}, err error)
	}
	sqlBuilder struct {
		buf *bytes.Buffer
		// True if the sql should not be interpolated
		isPrepared bool
		// Current Number of arguments, used by adapters that need positional placeholders
		currentArgPosition int
		args               []interface{}
		err                error
	}
)

func NewSQLBuilder(isPrepared bool) SQLBuilder {
	return &sqlBuilder{
		buf:                &bytes.Buffer{},
		isPrepared:         isPrepared,
		args:               make([]interface{}, 0),
		currentArgPosition: 1,
	}
}

func (b *sqlBuilder) Error() error {
	return b.err
}

func (b *sqlBuilder) SetError(err error) SQLBuilder {
	if b.err == nil {
		b.err = err
	}
	return b
}

func (b *sqlBuilder) Write(bs []byte) SQLBuilder {
	if b.err == nil {
		b.buf.Write(bs)
	}
	return b
}

func (b *sqlBuilder) WriteStrings(ss ...string) SQLBuilder {
	if b.err == nil {
		for _, s := range ss {
			b.buf.WriteString(s)
		}
	}
	return b
}

func (b *sqlBuilder) WriteRunes(rs ...rune) SQLBuilder {
	if b.err == nil {
		for _, r := range rs {
			b.buf.WriteRune(r)
		}
	}
	return b
}

// Returns true if the sql is a prepared statement
func (b *sqlBuilder) IsPrepared() bool {
	return b.isPrepared
}

// Returns true if the sql is a prepared statement
func (b *sqlBuilder) CurrentArgPosition() int {
	return b.currentArgPosition
}

// Adds an argument to the builder, used when IsPrepared is false
func (b *sqlBuilder) WriteArg(i ...interface{}) SQLBuilder {
	if b.err == nil {
		b.currentArgPosition += len(i)
		b.args = append(b.args, i...)
	}
	return b
}

// Returns the sql string, and arguments.
func (b *sqlBuilder) ToSQL() (sql string, args []interface{}, err error) {
	if b.err != nil {
		return sql, args, b.err
	}
	return b.buf.String(), b.args, nil
}
