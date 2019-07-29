package config

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/spf13/viper"
)

const Stash = "stash"
const Cache = "cache"
const Generated = "generated"
const Metadata = "metadata"
const Downloads = "downloads"
const Username = "username"
const Password = "password"

const Database = "database"

const Host = "host"
const Port = "port"

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

func SetPassword(value string) {
	// if blank, don't bother hashing; we want it to be blank
	if value == "" {
		Set(Password, "")
	} else {
		Set(Password, hashPassword(value))
	}
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

func GetUsername() string {
	return viper.GetString(Username)
}

func GetPasswordHash() string {
	return viper.GetString(Password)
}

func GetCredentials() (string, string) {
	if HasCredentials() {
		return viper.GetString(Username), viper.GetString(Password)
	}

	return "", ""
}

func HasCredentials() bool {
	if !viper.IsSet(Username) || !viper.IsSet(Password) {
		return false
	}

	username := GetUsername()
	pwHash := GetPasswordHash()

	return username != "" && pwHash != ""
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hash)
}

func ValidateCredentials(username string, password string) bool {
	if !HasCredentials() {
		// don't need to authenticate if no credentials saved
		return true
	}

	authUser, authPWHash := GetCredentials()

	err := bcrypt.CompareHashAndPassword([]byte(authPWHash), []byte(password))

	return username == authUser && err == nil
}

func IsValid() bool {
	setPaths := viper.IsSet(Stash) && viper.IsSet(Cache) && viper.IsSet(Generated) && viper.IsSet(Metadata)

	// TODO: check valid paths
	return setPaths
}
