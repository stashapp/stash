package goqu

var (
	// defaultPrepared is controlled by SetDefaultPrepared
	defaultPrepared bool
)

type prepared int

const (
	// zero value that defers to defaultPrepared
	preparedNoPreference prepared = iota

	// explicitly enabled via Prepared(true) on a dataset
	preparedEnabled

	// explicitly disabled via Prepared(false) on a dataset
	preparedDisabled
)

// Bool converts the ternary prepared state into a boolean. If the prepared
// state is preparedNoPreference, the value depends on the last value that
// SetDefaultPrepared was called with which is false by default.
func (p prepared) Bool() bool {
	if p == preparedNoPreference {
		return defaultPrepared
	} else if p == preparedEnabled {
		return true
	}

	return false
}

// preparedFromBool converts a bool from e.g. Prepared(true) into a prepared
// const.
func preparedFromBool(prepared bool) prepared {
	if prepared {
		return preparedEnabled
	}

	return preparedDisabled
}

// SetDefaultPrepared controls the default Prepared state of all datasets. If
// set to true, any new dataset will use prepared queries by default.
func SetDefaultPrepared(prepared bool) {
	defaultPrepared = prepared
}
