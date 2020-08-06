package config

import (
	"golang.org/x/crypto/bcrypt"

	"io/ioutil"
	"path/filepath"

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
const MaxSessionAge = "max_session_age"

const DefaultMaxSessionAge = 60 * 60 * 1 // 1 hours

const Database = "database"

const Exclude = "exclude"

// CalculateMD5 is the config key used to determine if MD5 should be calculated
// for video files.
const CalculateMD5 = "calculate_md5"

// VideoFileNamingAlgorithm is the config key used to determine what hash
// should be used when generating and using generated files for scenes.
const VideoFileNamingAlgorithm = "video_file_naming_algorithm"

const PreviewPreset = "preview_preset"

const MaxTranscodeSize = "max_transcode_size"
const MaxStreamingTranscodeSize = "max_streaming_transcode_size"

const PreviewSegmentDuration = "preview_segment_duration"
const previewSegmentDurationDefault = 0.75

const PreviewSegments = "preview_segments"
const previewSegmentsDefault = 12

const PreviewExcludeStart = "preview_exclude_start"
const previewExcludeStartDefault = "0"

const PreviewExcludeEnd = "preview_exclude_end"
const previewExcludeEndDefault = "0"

const Host = "host"
const Port = "port"
const ExternalHost = "external_host"

// key used to sign JWT tokens
const JWTSignKey = "jwt_secret_key"

// key used for session store
const SessionStoreKey = "session_store_key"

// scraping options
const ScrapersPath = "scrapers_path"
const ScraperUserAgent = "scraper_user_agent"
const ScraperCDPPath = "scraper_cdp_path"

// i18n
const Language = "language"

// served directories
// this should be manually configured only
const CustomServedFolders = "custom_served_folders"

// Interface options
const SoundOnPreview = "sound_on_preview"
const WallShowTitle = "wall_show_title"
const MaximumLoopDuration = "maximum_loop_duration"
const AutostartVideo = "autostart_video"
const ShowStudioAsText = "show_studio_as_text"
const CSSEnabled = "cssEnabled"
const WallPlayback = "wall_playback"

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

func GetJWTSignKey() []byte {
	return []byte(viper.GetString(JWTSignKey))
}

func GetSessionStoreKey() []byte {
	return []byte(viper.GetString(SessionStoreKey))
}

func GetDefaultScrapersPath() string {
	// default to the same directory as the config file
	configFileUsed := viper.ConfigFileUsed()
	configDir := filepath.Dir(configFileUsed)

	fn := filepath.Join(configDir, "scrapers")

	return fn
}

func GetExcludes() []string {
	return viper.GetStringSlice(Exclude)
}

func GetLanguage() string {
	ret := viper.GetString(Language)

	// default to English
	if ret == "" {
		return "en-US"
	}

	return ret
}

// IsCalculateMD5 returns true if MD5 checksums should be generated for
// scene video files.
func IsCalculateMD5() bool {
	return viper.GetBool(CalculateMD5)
}

// GetVideoFileNamingAlgorithm returns what hash algorithm should be used for
// naming generated scene video files.
func GetVideoFileNamingAlgorithm() models.HashAlgorithm {
	ret := viper.GetString(VideoFileNamingAlgorithm)

	// default to oshash
	if ret == "" {
		return models.HashAlgorithmOshash
	}

	return models.HashAlgorithm(ret)
}

func GetScrapersPath() string {
	return viper.GetString(ScrapersPath)
}

func GetScraperUserAgent() string {
	return viper.GetString(ScraperUserAgent)
}

// GetScraperCDPPath gets the path to the Chrome executable or remote address
// to an instance of Chrome.
func GetScraperCDPPath() string {
	return viper.GetString(ScraperCDPPath)
}

func GetHost() string {
	return viper.GetString(Host)
}

func GetPort() int {
	return viper.GetInt(Port)
}

func GetExternalHost() string {
	return viper.GetString(ExternalHost)
}

// GetPreviewSegmentDuration returns the duration of a single segment in a
// scene preview file, in seconds.
func GetPreviewSegmentDuration() float64 {
	return viper.GetFloat64(PreviewSegmentDuration)
}

// GetPreviewSegments returns the amount of segments in a scene preview file.
func GetPreviewSegments() int {
	return viper.GetInt(PreviewSegments)
}

// GetPreviewExcludeStart returns the configuration setting string for
// excluding the start of scene videos for preview generation. This can
// be in two possible formats. A float value is interpreted as the amount
// of seconds to exclude from the start of the video before it is included
// in the preview. If the value is suffixed with a '%' character (for example
// '2%'), then it is interpreted as a proportion of the total video duration.
func GetPreviewExcludeStart() string {
	return viper.GetString(PreviewExcludeStart)
}

// GetPreviewExcludeEnd returns the configuration setting string for
// excluding the end of scene videos for preview generation. A float value
// is interpreted as the amount of seconds to exclude from the end of the video
// when generating previews. If the value is suffixed with a '%' character,
// then it is interpreted as a proportion of the total video duration.
func GetPreviewExcludeEnd() string {
	return viper.GetString(PreviewExcludeEnd)
}

// GetPreviewPreset returns the preset when generating previews. Defaults to
// Slow.
func GetPreviewPreset() models.PreviewPreset {
	ret := viper.GetString(PreviewPreset)

	// default to slow
	if ret == "" {
		return models.PreviewPresetSlow
	}

	return models.PreviewPreset(ret)
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

// GetMaxSessionAge gets the maximum age for session cookies, in seconds.
// Session cookie expiry times are refreshed every request.
func GetMaxSessionAge() int {
	viper.SetDefault(MaxSessionAge, DefaultMaxSessionAge)
	return viper.GetInt(MaxSessionAge)
}

// GetCustomServedFolders gets the map of custom paths to their applicable
// filesystem locations
func GetCustomServedFolders() URLMap {
	return viper.GetStringMapString(CustomServedFolders)
}

// Interface options
func GetSoundOnPreview() bool {
	viper.SetDefault(SoundOnPreview, true)
	return viper.GetBool(SoundOnPreview)
}

func GetWallShowTitle() bool {
	viper.SetDefault(WallShowTitle, true)
	return viper.GetBool(WallShowTitle)
}

func GetWallPlayback() string {
	viper.SetDefault(WallPlayback, "video")
	return viper.GetString(WallPlayback)
}

func GetMaximumLoopDuration() int {
	viper.SetDefault(MaximumLoopDuration, 0)
	return viper.GetInt(MaximumLoopDuration)
}

func GetAutostartVideo() bool {
	viper.SetDefault(AutostartVideo, false)
	return viper.GetBool(AutostartVideo)
}

func GetShowStudioAsText() bool {
	viper.SetDefault(ShowStudioAsText, false)
	return viper.GetBool(ShowStudioAsText)
}

func GetCSSPath() string {
	// use custom.css in the same directory as the config file
	configFileUsed := viper.ConfigFileUsed()
	configDir := filepath.Dir(configFileUsed)

	fn := filepath.Join(configDir, "custom.css")

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
	if value != "Debug" && value != "Info" && value != "Warning" && value != "Error" && value != "Trace" {
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

func setDefaultValues() {
	viper.SetDefault(PreviewSegmentDuration, previewSegmentDurationDefault)
	viper.SetDefault(PreviewSegments, previewSegmentsDefault)
	viper.SetDefault(PreviewExcludeStart, previewExcludeStartDefault)
	viper.SetDefault(PreviewExcludeEnd, previewExcludeEndDefault)
}

// SetInitialConfig fills in missing required config fields
func SetInitialConfig() error {
	// generate some api keys
	const apiKeyLength = 32

	if string(GetJWTSignKey()) == "" {
		signKey := utils.GenerateRandomKey(apiKeyLength)
		Set(JWTSignKey, signKey)
	}

	if string(GetSessionStoreKey()) == "" {
		sessionStoreKey := utils.GenerateRandomKey(apiKeyLength)
		Set(SessionStoreKey, sessionStoreKey)
	}

	setDefaultValues()

	return Write()
}
