package config

import (
	"github.com/spf13/viper"
)

const Stash = "stash"
const Cache = "cache"
const Generated = "generated"
const Metadata = "metadata"
const Downloads = "downloads"

const Database = "database"

const Host = "host"
const Port = "port"

const Verbose = "verbose"
const VerboseLevel1 = 1
const VerboseLevel2 = 2

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

func Write() error {
	return viper.WriteConfig()
}

func GetStashPaths() []string {
	return viper.GetStringSlice(Stash)
}

func GetCachePath() string {
	return viper.GetString(Cache)
}

func GetGeneratedPath() string {
	return viper.GetString(Generated)
}

func GetMetadataPath() string {
	return viper.GetString(Metadata)
}

func GetDatabasePath() string {
	return viper.GetString(Database)
}

func GetHost() string {
	return viper.GetString(Host)
}

func GetPort() int {
	return viper.GetInt(Port)
}

func GetVerbose() int {
	return viper.GetInt(Verbose)
}

func IsValid() bool {
	setPaths := viper.IsSet(Stash) && viper.IsSet(Cache) && viper.IsSet(Generated) && viper.IsSet(Metadata)
	// TODO: check valid paths
	return setPaths
}
