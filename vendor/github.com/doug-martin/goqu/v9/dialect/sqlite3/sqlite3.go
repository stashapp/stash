package sqlite3

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func DialectOptions() *goqu.SQLDialectOptions {
	opts := goqu.DefaultDialectOptions()

	opts.SupportsReturn = false
	opts.SupportsOrderByOnUpdate = true
	opts.SupportsLimitOnUpdate = true
	opts.SupportsOrderByOnDelete = true
	opts.SupportsLimitOnDelete = true
	opts.SupportsConflictUpdateWhere = false
	opts.SupportsInsertIgnoreSyntax = true
	opts.SupportsConflictTarget = true
	opts.SupportsMultipleUpdateTables = false
	opts.WrapCompoundsInParens = false
	opts.SupportsDistinctOn = false
	opts.SupportsWindowFunction = false
	opts.SupportsLateral = false

	opts.PlaceHolderFragment = []byte("?")
	opts.IncludePlaceholderNum = false
	opts.QuoteRune = '`'
	opts.DefaultValuesFragment = []byte("")
	opts.True = []byte("1")
	opts.False = []byte("0")
	opts.TimeFormat = time.RFC3339Nano
	opts.BooleanOperatorLookup = map[exp.BooleanOperation][]byte{
		exp.EqOp:             []byte("="),
		exp.NeqOp:            []byte("!="),
		exp.GtOp:             []byte(">"),
		exp.GteOp:            []byte(">="),
		exp.LtOp:             []byte("<"),
		exp.LteOp:            []byte("<="),
		exp.InOp:             []byte("IN"),
		exp.NotInOp:          []byte("NOT IN"),
		exp.IsOp:             []byte("IS"),
		exp.IsNotOp:          []byte("IS NOT"),
		exp.LikeOp:           []byte("LIKE"),
		exp.NotLikeOp:        []byte("NOT LIKE"),
		exp.ILikeOp:          []byte("LIKE"),
		exp.NotILikeOp:       []byte("NOT LIKE"),
		exp.RegexpLikeOp:     []byte("REGEXP"),
		exp.RegexpNotLikeOp:  []byte("NOT REGEXP"),
		exp.RegexpILikeOp:    []byte("REGEXP"),
		exp.RegexpNotILikeOp: []byte("NOT REGEXP"),
	}
	opts.UseLiteralIsBools = false
	opts.BitwiseOperatorLookup = map[exp.BitwiseOperation][]byte{
		exp.BitwiseOrOp:         []byte("|"),
		exp.BitwiseAndOp:        []byte("&"),
		exp.BitwiseLeftShiftOp:  []byte("<<"),
		exp.BitwiseRightShiftOp: []byte(">>"),
	}
	opts.EscapedRunes = map[rune][]byte{
		'\'': []byte("''"),
	}
	opts.InsertIgnoreClause = []byte("INSERT OR IGNORE INTO ")
	opts.ConflictFragment = []byte(" ON CONFLICT ")
	opts.ConflictDoUpdateFragment = []byte(" DO UPDATE SET ")
	opts.ConflictDoNothingFragment = []byte(" DO NOTHING ")
	opts.ForUpdateFragment = []byte("")
	opts.OfFragment = []byte("")
	opts.NowaitFragment = []byte("")
	return opts
}

func init() {
	goqu.RegisterDialect("sqlite3", DialectOptions())
}
