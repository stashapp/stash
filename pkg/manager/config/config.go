package config

import (
	"golang.org/x/crypto/bcrypt"

	"io/ioutil"

	"github.com/spf13/viper"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

const Stash = "stash"
const Cache = "cache"
const Generated = "generated"
const Metadata = "metadata"
const Downloads = "downloads"
const Username = "username"
const Password = "password"

const Database = "database"

const MaxTranscodeSize = "max_transcode_size"
const MaxStreamingTranscodeSize = "max_streaming_transcode_size"

const Host = "host"
const Port = "port"

const CSSEnabled = "cssEnabled"

// Logging options
const LogFile = "logFile"
const LogOut = "logOut"
const LogLevel = "logLevel"
const LogAccess = "logAccess"

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

func GetMaxTranscodeSize() models.StreamingResolutionEnum {
	ret := viper.GetString(MaxTranscodeSize)

	// default to original
	if ret == "" {
		return models.StreamingResolutionEnumOriginal
	}

	return models.StreamingResolutionEnum(ret)
}

func GetMaxStreamingTranscodeSize() models.StreamingResolutionEnum {
	ret := viper.GetString(MaxStreamingTranscodeSize)

	// default to original
	if ret == "" {
		return models.StreamingResolutionEnumOriginal
	}

	return models.StreamingResolutionEnum(ret)
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

func GetCSSPath() string {
	// search for custom.css in current directory, then $HOME/.stash
	fn := "custom.css"
	exists, _ := utils.FileExists(fn)
	if !exists {
		fn = "$HOME/.stash/" + fn
	}

	return fn
}

func GetCSS() string {
	fn := GetCSSPath()

	exists, _ := utils.FileExists(fn)
	if !exists {
		return ""
	}

	buf, err := ioutil.ReadFile(fn)

	if err != nil {
		return ""
	}

	return string(buf)
}

func SetCSS(css string) {
	fn := GetCSSPath()

	buf := []byte(css)

	ioutil.WriteFile(fn, buf, 0777)
}

func GetCSSEnabled() bool {
	return viper.GetBool(CSSEnabled)
}

// GetLogFile returns the filename of the file to output logs to.
// An empty string means that file logging will be disabled.
func GetLogFile() string {
	return viper.GetString(LogFile)
}

// GetLogOut returns true if logging should be output to the terminal
// in addition to writing to a log file. Logging will be output to the
// terminal if file logging is disabled. Defaults to true.
func GetLogOut() bool {
	ret := true
	if viper.IsSet(LogOut) {
		ret = viper.GetBool(LogOut)
	}

	return ret
}

// GetLogLevel returns the lowest log level to write to the log.
// Should be one of "Debug", "Info", "Warning", "Error"
func GetLogLevel() string {
	const defaultValue = "Info"

	value := viper.GetString(LogLevel)
	if value != "Debug" && value != "Info" && value != "Warning" && value != "Error" {
		value = defaultValue
	}

	return value
}

// GetLogAccess returns true if http requests should be logged to the terminal.
// HTTP requests are not logged to the log file. Defaults to true.
func GetLogAccess() bool {
	ret := true
	if viper.IsSet(LogAccess) {
		ret = viper.GetBool(LogAccess)
	}

	return ret
}

func IsValid() bool {
	setPaths := viper.IsSet(Stash) && viper.IsSet(Cache) && viper.IsSet(Generated) && viper.IsSet(Metadata)

	// TODO: check valid paths
	return setPaths
}
