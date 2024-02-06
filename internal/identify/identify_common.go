package identify

import (
	"errors"
	"fmt"
)

var (
	ErrSkipSingleNamePerformer = errors.New("a performer was skipped because they only had a single name and no disambiguation")
)

func (e *MultipleMatchesFoundError) Error() string {
	return fmt.Sprintf("multiple matches found for %s", e.Source.Name)
}

func getFieldStrategy(strategy *FieldOptions) FieldStrategy {
	// if unset then default to MERGE
	fs := FieldStrategyMerge

	if strategy != nil && strategy.Strategy.IsValid() {
		fs = strategy.Strategy
	}

	return fs
}

func shouldSetSingleValueField(strategy *FieldOptions, hasExistingValue bool) bool {
	// if unset then default to MERGE
	fs := getFieldStrategy(strategy)

	if fs == FieldStrategyIgnore {
		return false
	}

	return !hasExistingValue || fs == FieldStrategyOverwrite
}
