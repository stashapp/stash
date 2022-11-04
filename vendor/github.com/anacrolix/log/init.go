package log

import (
	"os"
)

func init() {
	var err error
	rules, err = parseEnvRules()
	if err != nil {
		panic(err)
	}
	Default = loggerCore{
		nonZero:     true,
		filterLevel: Error,
		Handlers:    []Handler{DefaultHandler},
	}.withFilterLevelFromRules()
	Default.defaultLevel, _, err = levelFromString(os.Getenv("GO_LOG_DEFAULT_LEVEL"))
	if err != nil {
		panic(err)
	}
}
