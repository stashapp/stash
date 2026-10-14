package sqlite

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func pathHasWildcard(path string) bool {
	return strings.Contains(path, "*")
}

func pathToLikePattern(path string) string {
	var b strings.Builder
	b.Grow(len(path))

	for _, r := range path {
		switch r {
		case '\\', '%', '_':
			b.WriteRune('\\')
		case '*':
			b.WriteRune('%')
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}

func pathLike(col exp.IdentifierExpression, pattern string) exp.Expression {
	return goqu.L(`? LIKE ? ESCAPE '\'`, col, pathToLikePattern(pattern))
}

func pathEqNoCase(col exp.IdentifierExpression, value string) exp.Expression {
	return goqu.L("? = ? COLLATE NOCASE", col, value)
}
