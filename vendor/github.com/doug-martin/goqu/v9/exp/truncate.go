package exp

// Options to use when generating a TRUNCATE statement
type TruncateOptions struct {
	// Set to true to add CASCADE to the TRUNCATE statement
	Cascade bool
	// Set to true to add RESTRICT to the TRUNCATE statement
	Restrict bool
	// Set to true to specify IDENTITY options, (e.g. RESTART, CONTINUE) to the TRUNCATE statement
	Identity string
}
