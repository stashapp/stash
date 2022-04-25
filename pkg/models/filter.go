package models

import (
	"fmt"
	"io"
	"strconv"
)

type CriterionModifier string

const (
	// =
	CriterionModifierEquals CriterionModifier = "EQUALS"
	// !=
	CriterionModifierNotEquals CriterionModifier = "NOT_EQUALS"
	// >
	CriterionModifierGreaterThan CriterionModifier = "GREATER_THAN"
	// <
	CriterionModifierLessThan CriterionModifier = "LESS_THAN"
	// IS NULL
	CriterionModifierIsNull CriterionModifier = "IS_NULL"
	// IS NOT NULL
	CriterionModifierNotNull CriterionModifier = "NOT_NULL"
	// INCLUDES ALL
	CriterionModifierIncludesAll CriterionModifier = "INCLUDES_ALL"
	CriterionModifierIncludes    CriterionModifier = "INCLUDES"
	CriterionModifierExcludes    CriterionModifier = "EXCLUDES"
	// MATCHES REGEX
	CriterionModifierMatchesRegex CriterionModifier = "MATCHES_REGEX"
	// NOT MATCHES REGEX
	CriterionModifierNotMatchesRegex CriterionModifier = "NOT_MATCHES_REGEX"
	// >= AND <=
	CriterionModifierBetween CriterionModifier = "BETWEEN"
	// < OR >
	CriterionModifierNotBetween CriterionModifier = "NOT_BETWEEN"
)

var AllCriterionModifier = []CriterionModifier{
	CriterionModifierEquals,
	CriterionModifierNotEquals,
	CriterionModifierGreaterThan,
	CriterionModifierLessThan,
	CriterionModifierIsNull,
	CriterionModifierNotNull,
	CriterionModifierIncludesAll,
	CriterionModifierIncludes,
	CriterionModifierExcludes,
	CriterionModifierMatchesRegex,
	CriterionModifierNotMatchesRegex,
	CriterionModifierBetween,
	CriterionModifierNotBetween,
}

func (e CriterionModifier) IsValid() bool {
	switch e {
	case CriterionModifierEquals, CriterionModifierNotEquals, CriterionModifierGreaterThan, CriterionModifierLessThan, CriterionModifierIsNull, CriterionModifierNotNull, CriterionModifierIncludesAll, CriterionModifierIncludes, CriterionModifierExcludes, CriterionModifierMatchesRegex, CriterionModifierNotMatchesRegex, CriterionModifierBetween, CriterionModifierNotBetween:
		return true
	}
	return false
}

func (e CriterionModifier) String() string {
	return string(e)
}

func (e *CriterionModifier) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CriterionModifier(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CriterionModifier", str)
	}
	return nil
}

func (e CriterionModifier) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type StringCriterionInput struct {
	Value    string            `json:"value"`
	Modifier CriterionModifier `json:"modifier"`
}

type IntCriterionInput struct {
	Value    int               `json:"value"`
	Value2   *int              `json:"value2"`
	Modifier CriterionModifier `json:"modifier"`
}

type ResolutionCriterionInput struct {
	Value    ResolutionEnum    `json:"value"`
	Modifier CriterionModifier `json:"modifier"`
}

type HierarchicalMultiCriterionInput struct {
	Value    []string          `json:"value"`
	Modifier CriterionModifier `json:"modifier"`
	Depth    *int              `json:"depth"`
}

type MultiCriterionInput struct {
	Value    []string          `json:"value"`
	Modifier CriterionModifier `json:"modifier"`
}
