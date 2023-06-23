# casefolded [![PkgGoDev](https://pkg.go.dev/badge/github.com/fvbommel/sortorder/casefolded)](https://pkg.go.dev/github.com/fvbommel/sortorder/casefolded)

    import "github.com/fvbommel/sortorder/casefolded"

Case-folded sort orders and comparison functions.

These sort characters as the lowest unicode value that is equivalent to that character, ignoring case.

Not all Unicode special cases are supported.

This is a separate sub-package because this needs to pull in the Unicode tables in the standard library,
which can add significantly to the size of binaries.
