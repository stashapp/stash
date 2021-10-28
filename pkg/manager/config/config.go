package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"sync"
	// "github.com/sasha-s/go-deadlock" // if you have deadlock issues

	"golang.org/x/crypto/bcrypt"

	"github.com/spf13/viper"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

const Stash = "stash"
const Cache = "cache"
const Generated = "generated"
const Metadata = "metadata"
const Downloads = "downloads"
const ApiKey = "api_key"
const Username = "username"
const Password = "password"
const MaxSessionAge = "max_session_age"

const DefaultMaxSessionAge = 60 * 60 * 1 // 1 hours

const Database = "database"

const Exclude = "exclude"
const ImageExclude = "image_exclude"

const VideoExtensions = "video_extensions"

var defaultVideoExtensions = []string{"m4v", "mp4", "mov", "wmv", "avi", "mpg", "mpeg", "rmvb", "rm", "flv", "asf", "mkv", "webm"}

const ImageExtensions = "image_extensions"

var defaultImageExtensions = []string{"png", "jpg", "jpeg", "gif", "webp"}

const GalleryExtensions = "gallery_extensions"

var defaultGalleryExtensions = []string{"zip", "cbz"}

const CreateGalleriesFromFolders = "create_galleries_from_folders"

// CalculateMD5 is the config key used to determine if MD5 should be calculated
// for video files.
const CalculateMD5 = "calculate_md5"

// VideoFileNamingAlgorithm is the config key used to determine what hash
// should be used when generating and using generated files for scenes.
const VideoFileNamingAlgorithm = "video_file_naming_algorithm"

const MaxTranscodeSize = "max_transcode_size"
const MaxStreamingTranscodeSize = "max_streaming_transcode_size"

const ParallelTasks = "parallel_tasks"
const parallelTasksDefault = 1

const PreviewPreset = "preview_preset"

const PreviewAudio = "preview_audio"
const previewAudioDefault = true

const PreviewSegmentDuration = "preview_segment_duration"
const previewSegmentDurationDefault = 0.75

const PreviewSegments = "preview_segments"
const previewSegmentsDefault = 12

const PreviewExcludeStart = "preview_exclude_start"
const previewExcludeStartDefault = "0"

const PreviewExcludeEnd = "preview_exclude_end"
const previewExcludeEndDefault = "0"

const WriteImageThumbnails = "write_image_thumbnails"
const writeImageThumbnailsDefault = true

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
const ScraperCertCheck = "scraper_cert_check"
const ScraperCDPPath = "scraper_cdp_path"
const ScraperExcludeTagPatterns = "scraper_exclude_tag_patterns"

// stash-box options
const StashBoxes = "stash_boxes"

// plugin options
const PluginsPath = "plugins_path"

// i18n
const Language = "language"

// served directories
// this should be manually configured only
const CustomServedFolders = "custom_served_folders"

// UI directory. Overrides to serve the UI from a specific location
// rather than use the embedded UI.
const CustomUILocation = "custom_ui_location"

// Interface options
const MenuItems = "menu_items"

var defaultMenuItems = []string{"scenes", "images", "movies", "markers", "galleries", "performers", "studios", "tags"}

const SoundOnPreview = "sound_on_preview"
const WallShowTitle = "wall_show_title"
const CustomPerformerImageLocation = "custom_performer_image_location"
const MaximumLoopDuration = "maximum_loop_duration"
const AutostartVideo = "autostart_video"
const ShowStudioAsText = "show_studio_as_text"
const CSSEnabled = "cssEnabled"
const WallPlayback = "wall_playback"
const SlideshowDelay = "slideshow_delay"

const (
	DisableDropdownCreatePerformer = "disable_dropdown_create.performer"
	DisableDropdownCreateStudio    = "disable_dropdown_create.studio"
	DisableDropdownCreateTag       = "disable_dropdown_create.tag"
)

const HandyKey = "handy_key"
const FunscriptOffset = "funscript_offset"

// Default settings
const (
	DefaultIdentifySettings = "defaults.identify_task"

	DeleteFileDefault      = "defaults.delete_file"
	DeleteGeneratedDefault = "defaults.delete_generated"
)

// Security
const TrustedProxies = "trusted_proxies"
const dangerousAllowPublicWithoutAuth = "dangerous_allow_public_without_auth"
const dangerousAllowPublicWithoutAuthDefault = "false"
const SecurityTripwireAccessedFromPublicInternet = "security_tripwire_accessed_from_public_internet"
const securityTripwireAccessedFromPublicInternetDefault = ""

// DLNA options
const DLNAServerName = "dlna.server_name"
const DLNADefaultEnabled = "dlna.default_enabled"
const DLNADefaultIPWhitelist = "dlna.default_whitelist"
const DLNAInterfaces = "dlna.interfaces"

// Logging options
const LogFile = "logFile"
const LogOut = "logOut"
const LogLevel = "logLevel"
const LogAccess = "logAccess"

// File upload options
const MaxUploadSize = "max_upload_size"

type MissingConfigError struct {
	missingFields []string
}

func (e MissingConfigError) Error() string {
	return fmt.Sprintf("missing the following mandatory settings: %s", strings.Join(e.missingFields, ", "))
}

// StashBoxError represents configuration errors of Stash-Box
type StashBoxError struct {
	msg string
}

func (s *StashBoxError) Error() string {
	// "Stash-box" is a proper noun and is therefore capitcalized
	return "Stash-box: " + s.msg
}

type Instance struct {
	cpuProfilePath string
	isNewSystem    bool
	certFile       string
	keyFile        string
	sync.RWMutex
	// deadlock.RWMutex // for deadlock testing/issues
}

var instance *Instance

func GetInstance() *Instance {
	if instance == nil {
		instance = &Instance{}
	}
	return instance
}

func (i *Instance) IsNewSystem() bool {
	return i.isNewSystem
}

func (i *Instance) SetConfigFile(fn string) {
	i.Lock()
	defer i.Unlock()
	viper.SetConfigFile(fn)
}

func (i *Instance) InitTLS() {
	configDirectory := i.GetConfigPath()
	tlsPaths := []string{
		configDirectory,
		paths.GetStashHomeDirectory(),
	}

	i.certFile = utils.FindInPaths(tlsPaths, "stash.crt")
	i.keyFile = utils.FindInPaths(tlsPaths, "stash.key")
}

func (i *Instance) GetTLSFiles() (certFile, keyFile string) {
	return i.certFile, i.keyFile
}

func (i *Instance) HasTLSConfig() bool {
	certFile, keyFile := i.GetTLSFiles()
	return certFile != "" && keyFile != ""
}

// GetCPUProfilePath returns the path to the CPU profile file to output
// profiling info to. This is set only via a commandline flag. Returns an
// empty string if not set.
func (i *Instance) GetCPUProfilePath() string {
	return i.cpuProfilePath
}

func (i *Instance) Set(key string, value interface{}) {
	i.Lock()
	defer i.Unlock()
	viper.Set(key, value)
}

func (i *Instance) SetPassword(value string) {
	// if blank, don't bother hashing; we want it to be blank
	if value == "" {
		i.Set(Password, "")
	} else {
		i.Set(Password, hashPassword(value))
	}
}

func (i *Instance) Write() error {
	i.Lock()
	defer i.Unlock()
	return viper.WriteConfig()
}

// GetConfigFile returns the full path to the used configuration file.
func (i *Instance) GetConfigFile() string {
	i.RLock()
	defer i.RUnlock()
	return viper.ConfigFileUsed()
}

// GetConfigPath returns the path of the directory containing the used
// configuration file.
func (i *Instance) GetConfigPath() string {
	return filepath.Dir(i.GetConfigFile())
}

// GetDefaultDatabaseFilePath returns the default database filename,
// which is located in the same directory as the config file.
func (i *Instance) GetDefaultDatabaseFilePath() string {
	return filepath.Join(i.GetConfigPath(), "stash-go.sqlite")
}

func (i *Instance) GetStashPaths() []*models.StashConfig {
	i.RLock()
	defer i.RUnlock()
	var ret []*models.StashConfig
	if err := viper.UnmarshalKey(Stash, &ret); err != nil || len(ret) == 0 {
		// fallback to legacy format
		ss := viper.GetStringSlice(Stash)
		ret = nil
		for _, path := range ss {
			toAdd := &models.StashConfig{
				Path: path,
			}
			ret = append(ret, toAdd)
		}
	}

	return ret
}

func (i *Instance) GetConfigFilePath() string {
	i.RLock()
	defer i.RUnlock()
	return viper.ConfigFileUsed()
}

func (i *Instance) GetCachePath() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(Cache)
}

func (i *Instance) GetGeneratedPath() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(Generated)
}

func (i *Instance) GetMetadataPath() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(Metadata)
}

func (i *Instance) GetDatabasePath() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(Database)
}

func (i *Instance) GetJWTSignKey() []byte {
	i.RLock()
	defer i.RUnlock()
	return []byte(viper.GetString(JWTSignKey))
}

func (i *Instance) GetSessionStoreKey() []byte {
	i.RLock()
	defer i.RUnlock()
	return []byte(viper.GetString(SessionStoreKey))
}

func (i *Instance) GetDefaultScrapersPath() string {
	// default to the same directory as the config file

	fn := filepath.Join(i.GetConfigPath(), "scrapers")

	return fn
}

func (i *Instance) GetExcludes() []string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetStringSlice(Exclude)
}

func (i *Instance) GetImageExcludes() []string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetStringSlice(ImageExclude)
}

func (i *Instance) GetVideoExtensions() []string {
	i.RLock()
	defer i.RUnlock()
	ret := viper.GetStringSlice(VideoExtensions)
	if ret == nil {
		ret = defaultVideoExtensions
	}
	return ret
}

func (i *Instance) GetImageExtensions() []string {
	i.RLock()
	defer i.RUnlock()
	ret := viper.GetStringSlice(ImageExtensions)
	if ret == nil {
		ret = defaultImageExtensions
	}
	return ret
}

func (i *Instance) GetGalleryExtensions() []string {
	i.RLock()
	defer i.RUnlock()
	ret := viper.GetStringSlice(GalleryExtensions)
	if ret == nil {
		ret = defaultGalleryExtensions
	}
	return ret
}

func (i *Instance) GetCreateGalleriesFromFolders() bool {
	i.RLock()
	defer i.RUnlock()
	return viper.GetBool(CreateGalleriesFromFolders)
}

func (i *Instance) GetLanguage() string {
	i.RLock()
	defer i.RUnlock()
	ret := viper.GetString(Language)

	// default to English
	if ret == "" {
		return "en-US"
	}

	return ret
}

// IsCalculateMD5 returns true if MD5 checksums should be generated for
// scene video files.
func (i *Instance) IsCalculateMD5() bool {
	i.RLock()
	defer i.RUnlock()
	return viper.GetBool(CalculateMD5)
}

// GetVideoFileNamingAlgorithm returns what hash algorithm should be used for
// naming generated scene video files.
func (i *Instance) GetVideoFileNamingAlgorithm() models.HashAlgorithm {
	i.RLock()
	defer i.RUnlock()
	ret := viper.GetString(VideoFileNamingAlgorithm)

	// default to oshash
	if ret == "" {
		return models.HashAlgorithmOshash
	}

	return models.HashAlgorithm(ret)
}

func (i *Instance) GetScrapersPath() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(ScrapersPath)
}

func (i *Instance) GetScraperUserAgent() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(ScraperUserAgent)
}

// GetScraperCDPPath gets the path to the Chrome executable or remote address
// to an instance of Chrome.
func (i *Instance) GetScraperCDPPath() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(ScraperCDPPath)
}

// GetScraperCertCheck returns true if the scraper should check for insecure
// certificates when fetching an image or a page.
func (i *Instance) GetScraperCertCheck() bool {
	i.RLock()
	defer i.RUnlock()
	ret := true
	if viper.IsSet(ScraperCertCheck) {
		ret = viper.GetBool(ScraperCertCheck)
	}

	return ret
}

func (i *Instance) GetScraperExcludeTagPatterns() []string {
	i.RLock()
	defer i.RUnlock()
	var ret []string
	if viper.IsSet(ScraperExcludeTagPatterns) {
		ret = viper.GetStringSlice(ScraperExcludeTagPatterns)
	}

	return ret
}

func (i *Instance) GetStashBoxes() models.StashBoxes {
	i.RLock()
	defer i.RUnlock()
	var boxes models.StashBoxes
	if err := viper.UnmarshalKey(StashBoxes, &boxes); err != nil {
		logger.Warnf("error in unmarshalkey: %v", err)
	}

	return boxes
}

func (i *Instance) GetDefaultPluginsPath() string {
	// default to the same directory as the config file
	fn := filepath.Join(i.GetConfigPath(), "plugins")

	return fn
}

func (i *Instance) GetPluginsPath() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(PluginsPath)
}

func (i *Instance) GetHost() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(Host)
}

func (i *Instance) GetPort() int {
	i.RLock()
	defer i.RUnlock()
	return viper.GetInt(Port)
}

func (i *Instance) GetExternalHost() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(ExternalHost)
}

// GetPreviewSegmentDuration returns the duration of a single segment in a
// scene preview file, in seconds.
func (i *Instance) GetPreviewSegmentDuration() float64 {
	i.RLock()
	defer i.RUnlock()
	return viper.GetFloat64(PreviewSegmentDuration)
}

// GetParallelTasks returns the number of parallel tasks that should be started
// by scan or generate task.
func (i *Instance) GetParallelTasks() int {
	i.RLock()
	defer i.RUnlock()
	return viper.GetInt(ParallelTasks)
}

func (i *Instance) GetParallelTasksWithAutoDetection() int {
	i.RLock()
	defer i.RUnlock()
	parallelTasks := viper.GetInt(ParallelTasks)
	if parallelTasks <= 0 {
		parallelTasks = (runtime.NumCPU() / 4) + 1
	}
	return parallelTasks
}

func (i *Instance) GetPreviewAudio() bool {
	i.RLock()
	defer i.RUnlock()
	return viper.GetBool(PreviewAudio)
}

// GetPreviewSegments returns the amount of segments in a scene preview file.
func (i *Instance) GetPreviewSegments() int {
	i.RLock()
	defer i.RUnlock()
	return viper.GetInt(PreviewSegments)
}

// GetPreviewExcludeStart returns the configuration setting string for
// excluding the start of scene videos for preview generation. This can
// be in two possible formats. A float value is interpreted as the amount
// of seconds to exclude from the start of the video before it is included
// in the preview. If the value is suffixed with a '%' character (for example
// '2%'), then it is interpreted as a proportion of the total video duration.
func (i *Instance) GetPreviewExcludeStart() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(PreviewExcludeStart)
}

// GetPreviewExcludeEnd returns the configuration setting string for
// excluding the end of scene videos for preview generation. A float value
// is interpreted as the amount of seconds to exclude from the end of the video
// when generating previews. If the value is suffixed with a '%' character,
// then it is interpreted as a proportion of the total video duration.
func (i *Instance) GetPreviewExcludeEnd() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(PreviewExcludeEnd)
}

// GetPreviewPreset returns the preset when generating previews. Defaults to
// Slow.
func (i *Instance) GetPreviewPreset() models.PreviewPreset {
	i.RLock()
	defer i.RUnlock()
	ret := viper.GetString(PreviewPreset)

	// default to slow
	if ret == "" {
		return models.PreviewPresetSlow
	}

	return models.PreviewPreset(ret)
}

func (i *Instance) GetMaxTranscodeSize() models.StreamingResolutionEnum {
	i.RLock()
	defer i.RUnlock()
	ret := viper.GetString(MaxTranscodeSize)

	// default to original
	if ret == "" {
		return models.StreamingResolutionEnumOriginal
	}

	return models.StreamingResolutionEnum(ret)
}

func (i *Instance) GetMaxStreamingTranscodeSize() models.StreamingResolutionEnum {
	i.RLock()
	defer i.RUnlock()
	ret := viper.GetString(MaxStreamingTranscodeSize)

	// default to original
	if ret == "" {
		return models.StreamingResolutionEnumOriginal
	}

	return models.StreamingResolutionEnum(ret)
}

// IsWriteImageThumbnails returns true if image thumbnails should be written
// to disk after generating on the fly.
func (i *Instance) IsWriteImageThumbnails() bool {
	i.RLock()
	defer i.RUnlock()
	return viper.GetBool(WriteImageThumbnails)
}

func (i *Instance) GetAPIKey() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(ApiKey)
}

func (i *Instance) GetUsername() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(Username)
}

func (i *Instance) GetPasswordHash() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(Password)
}

func (i *Instance) GetCredentials() (string, string) {
	if i.HasCredentials() {
		i.RLock()
		defer i.RUnlock()
		return viper.GetString(Username), viper.GetString(Password)
	}

	return "", ""
}

func (i *Instance) HasCredentials() bool {
	i.RLock()
	defer i.RUnlock()
	if !viper.IsSet(Username) || !viper.IsSet(Password) {
		return false
	}

	username := viper.GetString(Username)
	pwHash := viper.GetString(Password)

	return username != "" && pwHash != ""
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hash)
}

func (i *Instance) ValidateCredentials(username string, password string) bool {
	if !i.HasCredentials() {
		// don't need to authenticate if no credentials saved
		return true
	}

	authUser, authPWHash := i.GetCredentials()

	err := bcrypt.CompareHashAndPassword([]byte(authPWHash), []byte(password))

	return username == authUser && err == nil
}

var stashBoxRe = regexp.MustCompile("^http.*graphql$")

func (i *Instance) ValidateStashBoxes(boxes []*models.StashBoxInput) error {
	isMulti := len(boxes) > 1

	for _, box := range boxes {
		// Validate each stash-box configuration field, return on error
		if box.APIKey == "" {
			return &StashBoxError{msg: "API Key cannot be blank"}
		}

		if box.Endpoint == "" {
			return &StashBoxError{msg: "endpoint cannot be blank"}
		}

		if !stashBoxRe.Match([]byte(box.Endpoint)) {
			return &StashBoxError{msg: "endpoint is invalid"}
		}

		if isMulti && box.Name == "" {
			return &StashBoxError{msg: "name cannot be blank"}
		}
	}

	return nil
}

// GetMaxSessionAge gets the maximum age for session cookies, in seconds.
// Session cookie expiry times are refreshed every request.
func (i *Instance) GetMaxSessionAge() int {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(MaxSessionAge, DefaultMaxSessionAge)
	return viper.GetInt(MaxSessionAge)
}

// GetCustomServedFolders gets the map of custom paths to their applicable
// filesystem locations
func (i *Instance) GetCustomServedFolders() URLMap {
	i.RLock()
	defer i.RUnlock()
	return viper.GetStringMapString(CustomServedFolders)
}

func (i *Instance) GetCustomUILocation() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(CustomUILocation)
}

// Interface options
func (i *Instance) GetMenuItems() []string {
	i.RLock()
	defer i.RUnlock()
	if viper.IsSet(MenuItems) {
		return viper.GetStringSlice(MenuItems)
	}
	return defaultMenuItems
}

func (i *Instance) GetSoundOnPreview() bool {
	i.RLock()
	defer i.RUnlock()
	return viper.GetBool(SoundOnPreview)
}

func (i *Instance) GetWallShowTitle() bool {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(WallShowTitle, true)
	return viper.GetBool(WallShowTitle)
}

func (i *Instance) GetCustomPerformerImageLocation() string {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(CustomPerformerImageLocation, "")
	return viper.GetString(CustomPerformerImageLocation)
}

func (i *Instance) GetWallPlayback() string {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(WallPlayback, "video")
	return viper.GetString(WallPlayback)
}

func (i *Instance) GetMaximumLoopDuration() int {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(MaximumLoopDuration, 0)
	return viper.GetInt(MaximumLoopDuration)
}

func (i *Instance) GetAutostartVideo() bool {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(AutostartVideo, false)
	return viper.GetBool(AutostartVideo)
}

func (i *Instance) GetShowStudioAsText() bool {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(ShowStudioAsText, false)
	return viper.GetBool(ShowStudioAsText)
}

func (i *Instance) GetSlideshowDelay() int {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(SlideshowDelay, 5000)
	return viper.GetInt(SlideshowDelay)
}

func (i *Instance) GetDisableDropdownCreate() *models.ConfigDisableDropdownCreate {
	i.Lock()
	defer i.Unlock()

	return &models.ConfigDisableDropdownCreate{
		Performer: viper.GetBool(DisableDropdownCreatePerformer),
		Studio:    viper.GetBool(DisableDropdownCreateStudio),
		Tag:       viper.GetBool(DisableDropdownCreateTag),
	}
}

func (i *Instance) GetCSSPath() string {
	i.RLock()
	defer i.RUnlock()
	// use custom.css in the same directory as the config file
	configFileUsed := viper.ConfigFileUsed()
	configDir := filepath.Dir(configFileUsed)

	fn := filepath.Join(configDir, "custom.css")

	return fn
}

func (i *Instance) GetCSS() string {
	fn := i.GetCSSPath()

	exists, _ := utils.FileExists(fn)
	if !exists {
		return ""
	}

	buf, err := os.ReadFile(fn)

	if err != nil {
		return ""
	}

	return string(buf)
}

func (i *Instance) SetCSS(css string) {
	i.RLock()
	defer i.RUnlock()
	fn := i.GetCSSPath()

	buf := []byte(css)

	if err := os.WriteFile(fn, buf, 0777); err != nil {
		logger.Warnf("error while writing %v bytes to %v: %v", len(buf), fn, err)
	}
}

func (i *Instance) GetCSSEnabled() bool {
	i.RLock()
	defer i.RUnlock()
	return viper.GetBool(CSSEnabled)
}

func (i *Instance) GetHandyKey() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(HandyKey)
}

func (i *Instance) GetFunscriptOffset() int {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(FunscriptOffset, 0)
	return viper.GetInt(FunscriptOffset)
}

func (i *Instance) GetDeleteFileDefault() bool {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(DeleteFileDefault, false)
	return viper.GetBool(DeleteFileDefault)
}

func (i *Instance) GetDeleteGeneratedDefault() bool {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(DeleteGeneratedDefault, true)
	return viper.GetBool(DeleteGeneratedDefault)
}

// GetDefaultIdentifySettings returns the default Identify task settings.
// Returns nil if the settings could not be unmarshalled, or if it
// has not been set.
func (i *Instance) GetDefaultIdentifySettings() *models.IdentifyMetadataTaskOptions {
	i.RLock()
	defer i.RUnlock()

	if viper.IsSet(DefaultIdentifySettings) {
		var ret models.IdentifyMetadataTaskOptions
		if err := viper.UnmarshalKey(DefaultIdentifySettings, &ret); err != nil {
			return nil
		}
		return &ret
	}

	return nil
}

// GetTrustedProxies returns a comma separated list of ip addresses that should allow proxying.
// When empty, allow from any private network
func (i *Instance) GetTrustedProxies() []string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetStringSlice(TrustedProxies)
}

// GetDangerousAllowPublicWithoutAuth determines if the security feature is enabled.
// See https://github.com/stashapp/stash/wiki/Authentication-Required-When-Accessing-Stash-From-the-Internet
func (i *Instance) GetDangerousAllowPublicWithoutAuth() bool {
	i.RLock()
	defer i.RUnlock()
	return viper.GetBool(dangerousAllowPublicWithoutAuth)
}

// GetSecurityTripwireAccessedFromPublicInternet returns a public IP address if stash
// has been accessed from the public internet, with no auth enabled, and
// DangerousAllowPublicWithoutAuth disabled. Returns an empty string otherwise.
func (i *Instance) GetSecurityTripwireAccessedFromPublicInternet() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(SecurityTripwireAccessedFromPublicInternet)
}

// GetDLNAServerName returns the visible name of the DLNA server. If empty,
// "stash" will be used.
func (i *Instance) GetDLNAServerName() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(DLNAServerName)
}

// GetDLNADefaultEnabled returns true if the DLNA is enabled by default.
func (i *Instance) GetDLNADefaultEnabled() bool {
	i.RLock()
	defer i.RUnlock()
	return viper.GetBool(DLNADefaultEnabled)
}

// GetDLNADefaultIPWhitelist returns a list of IP addresses/wildcards that
// are allowed to use the DLNA service.
func (i *Instance) GetDLNADefaultIPWhitelist() []string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetStringSlice(DLNADefaultIPWhitelist)
}

// GetDLNAInterfaces returns a list of interface names to expose DLNA on. If
// empty, runs on all interfaces.
func (i *Instance) GetDLNAInterfaces() []string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetStringSlice(DLNAInterfaces)
}

// GetLogFile returns the filename of the file to output logs to.
// An empty string means that file logging will be disabled.
func (i *Instance) GetLogFile() string {
	i.RLock()
	defer i.RUnlock()
	return viper.GetString(LogFile)
}

// GetLogOut returns true if logging should be output to the terminal
// in addition to writing to a log file. Logging will be output to the
// terminal if file logging is disabled. Defaults to true.
func (i *Instance) GetLogOut() bool {
	i.RLock()
	defer i.RUnlock()
	ret := true
	if viper.IsSet(LogOut) {
		ret = viper.GetBool(LogOut)
	}

	return ret
}

// GetLogLevel returns the lowest log level to write to the log.
// Should be one of "Debug", "Info", "Warning", "Error"
func (i *Instance) GetLogLevel() string {
	i.RLock()
	defer i.RUnlock()
	const defaultValue = "Info"

	value := viper.GetString(LogLevel)
	if value != "Debug" && value != "Info" && value != "Warning" && value != "Error" && value != "Trace" {
		value = defaultValue
	}

	return value
}

// GetLogAccess returns true if http requests should be logged to the terminal.
// HTTP requests are not logged to the log file. Defaults to true.
func (i *Instance) GetLogAccess() bool {
	i.RLock()
	defer i.RUnlock()
	ret := true
	if viper.IsSet(LogAccess) {
		ret = viper.GetBool(LogAccess)
	}

	return ret
}

// Max allowed graphql upload size in megabytes
func (i *Instance) GetMaxUploadSize() int64 {
	i.RLock()
	defer i.RUnlock()
	ret := int64(1024)
	if viper.IsSet(MaxUploadSize) {
		ret = viper.GetInt64(MaxUploadSize)
	}
	return ret << 20
}

// ActivatePublicAccessTripwire sets the security_tripwire_accessed_from_public_internet
// config field to the provided IP address to indicate that stash has been accessed
// from this public IP without authentication.
func (i *Instance) ActivatePublicAccessTripwire(requestIP string) error {
	i.Set(SecurityTripwireAccessedFromPublicInternet, requestIP)
	return i.Write()
}

func (i *Instance) Validate() error {
	i.RLock()
	defer i.RUnlock()
	mandatoryPaths := []string{
		Database,
		Generated,
	}

	var missingFields []string

	for _, p := range mandatoryPaths {
		if !viper.IsSet(p) || viper.GetString(p) == "" {
			missingFields = append(missingFields, p)
		}
	}

	if len(missingFields) > 0 {
		return MissingConfigError{
			missingFields: missingFields,
		}
	}

	return nil
}

func (i *Instance) SetChecksumDefaultValues(defaultAlgorithm models.HashAlgorithm, usingMD5 bool) {
	i.Lock()
	defer i.Unlock()
	viper.SetDefault(VideoFileNamingAlgorithm, defaultAlgorithm)
	viper.SetDefault(CalculateMD5, usingMD5)
}

func (i *Instance) setDefaultValues(write bool) error {
	// read data before write lock scope
	defaultDatabaseFilePath := i.GetDefaultDatabaseFilePath()
	defaultScrapersPath := i.GetDefaultScrapersPath()
	defaultPluginsPath := i.GetDefaultPluginsPath()

	i.Lock()
	defer i.Unlock()
	viper.SetDefault(ParallelTasks, parallelTasksDefault)
	viper.SetDefault(PreviewSegmentDuration, previewSegmentDurationDefault)
	viper.SetDefault(PreviewSegments, previewSegmentsDefault)
	viper.SetDefault(PreviewExcludeStart, previewExcludeStartDefault)
	viper.SetDefault(PreviewExcludeEnd, previewExcludeEndDefault)
	viper.SetDefault(PreviewAudio, previewAudioDefault)
	viper.SetDefault(SoundOnPreview, false)

	viper.SetDefault(WriteImageThumbnails, writeImageThumbnailsDefault)

	viper.SetDefault(Database, defaultDatabaseFilePath)

	viper.SetDefault(dangerousAllowPublicWithoutAuth, dangerousAllowPublicWithoutAuthDefault)
	viper.SetDefault(SecurityTripwireAccessedFromPublicInternet, securityTripwireAccessedFromPublicInternetDefault)

	// Set generated to the metadata path for backwards compat
	viper.SetDefault(Generated, viper.GetString(Metadata))

	// Set default scrapers and plugins paths
	viper.SetDefault(ScrapersPath, defaultScrapersPath)
	viper.SetDefault(PluginsPath, defaultPluginsPath)
	if write {
		return viper.WriteConfig()
	}

	return nil
}

// SetInitialConfig fills in missing required config fields
func (i *Instance) SetInitialConfig() error {
	return i.setInitialConfig(true)
}

// SetInitialMemoryConfig fills in missing required config fields without writing the configuration
func (i *Instance) SetInitialMemoryConfig() error {
	return i.setInitialConfig(false)
}

func (i *Instance) setInitialConfig(write bool) error {
	// generate some api keys
	const apiKeyLength = 32

	if string(i.GetJWTSignKey()) == "" {
		signKey := utils.GenerateRandomKey(apiKeyLength)
		i.Set(JWTSignKey, signKey)
	}

	if string(i.GetSessionStoreKey()) == "" {
		sessionStoreKey := utils.GenerateRandomKey(apiKeyLength)
		i.Set(SessionStoreKey, sessionStoreKey)
	}

	return i.setDefaultValues(write)
}

func (i *Instance) FinalizeSetup() {
	i.isNewSystem = false
}
