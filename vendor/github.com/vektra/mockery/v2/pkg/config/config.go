package config

// SemVer is the version of mockery at build time.
var SemVer = "0.0.0-dev"

type Config struct {
	All           bool
	BuildTags     string
	Case          string
	Config        string
	Cpuprofile    string
	Dir           string
	DryRun        bool `mapstructure:"dry-run"`
	FileName      string
	InPackage     bool
	KeepTree      bool
	LogLevel      string `mapstructure:"log-level"`
	Name          string
	Note          string
	Outpkg        string
	Packageprefix string
	Output        string
	Print         bool
	Profile       string
	Quiet         bool
	Recursive     bool
	SrcPkg        string
	// StructName overrides the name given to the mock struct and should only be nonempty
	// when generating for an exact match (non regex expression in -name).
	StructName string
	Tags       string
	TestOnly   bool
	Version    bool
}
