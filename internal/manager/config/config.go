package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"sync"
	// "github.com/sasha-s/go-deadlock" // if you have deadlock issues

	"golang.org/x/crypto/bcrypt"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"

	"github.com/stashapp/stash/internal/identify"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/hash"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
)

const (
	Stash               = "stash"
	Cache               = "cache"
	BackupDirectoryPath = "backup_directory_path"
	Generated           = "generated"
	Metadata            = "metadata"
	BlobsPath           = "blobs_path"
	Downloads           = "downloads"
	ApiKey              = "api_key"
	Username            = "username"
	Password            = "password"
	MaxSessionAge       = "max_session_age"

	FFMpegPath  = "ffmpeg_path"
	FFProbePath = "ffprobe_path"

	BlobsStorage = "blobs_storage"

	DefaultMaxSessionAge = 60 * 60 * 1 // 1 hours

	Database = "database"

	Exclude      = "exclude"
	ImageExclude = "image_exclude"

	VideoExtensions            = "video_extensions"
	ImageExtensions            = "image_extensions"
	GalleryExtensions          = "gallery_extensions"
	CreateGalleriesFromFolders = "create_galleries_from_folders"

	// CalculateMD5 is the config key used to determine if MD5 should be calculated
	// for video files.
	CalculateMD5 = "calculate_md5"

	// VideoFileNamingAlgorithm is the config key used to determine what hash
	// should be used when generating and using generated files for scenes.
	VideoFileNamingAlgorithm = "video_file_naming_algorithm"

	MaxTranscodeSize          = "max_transcode_size"
	MaxStreamingTranscodeSize = "max_streaming_transcode_size"

	// ffmpeg extra args options
	TranscodeInputArgs      = "ffmpeg.transcode.input_args"
	TranscodeOutputArgs     = "ffmpeg.transcode.output_args"
	LiveTranscodeInputArgs  = "ffmpeg.live_transcode.input_args"
	LiveTranscodeOutputArgs = "ffmpeg.live_transcode.output_args"

	ParallelTasks        = "parallel_tasks"
	parallelTasksDefault = 1

	PreviewPreset                 = "preview_preset"
	TranscodeHardwareAcceleration = "ffmpeg.hardware_acceleration"

	SequentialScanning        = "sequential_scanning"
	SequentialScanningDefault = false

	PreviewAudio        = "preview_audio"
	previewAudioDefault = true

	PreviewSegmentDuration        = "preview_segment_duration"
	previewSegmentDurationDefault = 0.75

	PreviewSegments        = "preview_segments"
	previewSegmentsDefault = 12

	PreviewExcludeStart        = "preview_exclude_start"
	previewExcludeStartDefault = "0"

	PreviewExcludeEnd        = "preview_exclude_end"
	previewExcludeEndDefault = "0"

	WriteImageThumbnails        = "write_image_thumbnails"
	writeImageThumbnailsDefault = true

	CreateImageClipsFromVideos        = "create_image_clip_from_videos"
	createImageClipsFromVideosDefault = false

	Host        = "host"
	hostDefault = "0.0.0.0"

	Port        = "port"
	portDefault = 9999

	ExternalHost = "external_host"

	// http proxy url if required
	Proxy = "proxy"

	// urls or IPs that should not use the proxy
	NoProxy        = "no_proxy"
	noProxyDefault = "localhost,127.0.0.1,192.168.0.0/16,10.0.0.0/8,172.16.0.0/12"

	// key used to sign JWT tokens
	JWTSignKey = "jwt_secret_key"

	// key used for session store
	SessionStoreKey = "session_store_key"

	// scraping options
	ScrapersPath              = "scrapers_path"
	ScraperUserAgent          = "scraper_user_agent"
	ScraperCertCheck          = "scraper_cert_check"
	ScraperCDPPath            = "scraper_cdp_path"
	ScraperExcludeTagPatterns = "scraper_exclude_tag_patterns"

	// stash-box options
	StashBoxes = "stash_boxes"

	PythonPath = "python_path"

	// plugin options
	PluginsPath          = "plugins_path"
	PluginsSetting       = "plugins.settings"
	PluginsSettingPrefix = PluginsSetting + "."
	DisabledPlugins      = "plugins.disabled"

	sourceDefaultPath = "community"
	sourceDefaultName = "Community (stable)"

	PluginPackageSources        = "plugins.package_sources"
	pluginPackageSourcesDefault = "https://stashapp.github.io/CommunityScripts/stable/index.yml"

	ScraperPackageSources        = "scrapers.package_sources"
	scraperPackageSourcesDefault = "https://stashapp.github.io/CommunityScrapers/stable/index.yml"

	// i18n
	Language = "language"

	// served directories
	// this should be manually configured only
	CustomServedFolders = "custom_served_folders"

	// UI directory. Overrides to serve the UI from a specific location
	// rather than use the embedded UI.
	UILocation = "ui_location"

	// backwards compatible name
	LegacyCustomUILocation = "custom_ui_location"

	// Gallery Cover Regex
	GalleryCoverRegex        = "gallery_cover_regex"
	galleryCoverRegexDefault = `(poster|cover|folder|board)\.[^\.]+$`

	// Interface options
	MenuItems = "menu_items"

	SoundOnPreview = "sound_on_preview"

	WallShowTitle        = "wall_show_title"
	defaultWallShowTitle = true

	CustomPerformerImageLocation        = "custom_performer_image_location"
	MaximumLoopDuration                 = "maximum_loop_duration"
	AutostartVideo                      = "autostart_video"
	AutostartVideoOnPlaySelected        = "autostart_video_on_play_selected"
	autostartVideoOnPlaySelectedDefault = true
	ContinuePlaylistDefault             = "continue_playlist_default"
	ShowStudioAsText                    = "show_studio_as_text"
	CSSEnabled                          = "cssenabled"
	JavascriptEnabled                   = "javascriptenabled"
	CustomLocalesEnabled                = "customlocalesenabled"

	ShowScrubber        = "show_scrubber"
	showScrubberDefault = true

	WallPlayback        = "wall_playback"
	defaultWallPlayback = "video"

	// Image lightbox options
	legacyImageLightboxSlideshowDelay       = "slideshow_delay"
	ImageLightboxSlideshowDelay             = "image_lightbox.slideshow_delay"
	ImageLightboxDisplayModeKey             = "image_lightbox.display_mode"
	ImageLightboxScaleUp                    = "image_lightbox.scale_up"
	ImageLightboxResetZoomOnNav             = "image_lightbox.reset_zoom_on_nav"
	ImageLightboxScrollModeKey              = "image_lightbox.scroll_mode"
	ImageLightboxScrollAttemptsBeforeChange = "image_lightbox.scroll_attempts_before_change"

	UI = "ui"

	defaultImageLightboxSlideshowDelay = 5

	DisableDropdownCreatePerformer = "disable_dropdown_create.performer"
	DisableDropdownCreateStudio    = "disable_dropdown_create.studio"
	DisableDropdownCreateTag       = "disable_dropdown_create.tag"
	DisableDropdownCreateMovie     = "disable_dropdown_create.movie"

	HandyKey                       = "handy_key"
	FunscriptOffset                = "funscript_offset"
	UseStashHostedFunscript        = "use_stash_hosted_funscript"
	useStashHostedFunscriptDefault = false

	DrawFunscriptHeatmapRange        = "draw_funscript_heatmap_range"
	drawFunscriptHeatmapRangeDefault = true

	ThemeColor        = "theme_color"
	DefaultThemeColor = "#202b33"

	// Security
	dangerousAllowPublicWithoutAuth                   = "dangerous_allow_public_without_auth"
	dangerousAllowPublicWithoutAuthDefault            = "false"
	SecurityTripwireAccessedFromPublicInternet        = "security_tripwire_accessed_from_public_internet"
	securityTripwireAccessedFromPublicInternetDefault = ""

	sslCertPath = "ssl_cert_path"
	sslKeyPath  = "ssl_key_path"

	// DLNA options
	DLNAServerName         = "dlna.server_name"
	DLNADefaultEnabled     = "dlna.default_enabled"
	DLNADefaultIPWhitelist = "dlna.default_whitelist"
	DLNAInterfaces         = "dlna.interfaces"

	DLNAVideoSortOrder        = "dlna.video_sort_order"
	dlnaVideoSortOrderDefault = "title"

	DLNAPort        = "dlna.port"
	DLNAPortDefault = 1338

	// Logging options
	LogFile          = "logfile"
	LogOut           = "logout"
	defaultLogOut    = true
	LogLevel         = "loglevel"
	defaultLogLevel  = "Info"
	LogAccess        = "logaccess"
	defaultLogAccess = true

	// Default settings
	DefaultScanSettings     = "defaults.scan_task"
	DefaultIdentifySettings = "defaults.identify_task"
	DefaultAutoTagSettings  = "defaults.auto_tag_task"
	DefaultGenerateSettings = "defaults.generate_task"

	DeleteFileDefault             = "defaults.delete_file"
	DeleteGeneratedDefault        = "defaults.delete_generated"
	deleteGeneratedDefaultDefault = true

	// Desktop Integration Options
	NoBrowser                           = "nobrowser"
	NoBrowserDefault                    = false
	NotificationsEnabled                = "notifications_enabled"
	NotificationsEnabledDefault         = true
	ShowOneTimeMovedNotification        = "show_one_time_moved_notification"
	ShowOneTimeMovedNotificationDefault = false

	// File upload options
	MaxUploadSize = "max_upload_size"

	// Developer options
	ExtraBlobsPaths = "developer_options.extra_blob_paths"
)

// slice default values
var (
	defaultVideoExtensions   = []string{"m4v", "mp4", "mov", "wmv", "avi", "mpg", "mpeg", "rmvb", "rm", "flv", "asf", "mkv", "webm"}
	defaultImageExtensions   = []string{"png", "jpg", "jpeg", "gif", "webp"}
	defaultGalleryExtensions = []string{"zip", "cbz"}
	defaultMenuItems         = []string{"scenes", "images", "movies", "markers", "galleries", "performers", "studios", "tags"}
)

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

type Config struct {
	// main instance - backed by config file
	main *koanf.Koanf

	// override instance - populated from flags/environment
	// not written to config file
	overrides *koanf.Koanf

	filePath    string
	isNewSystem bool
	// configUpdates  chan int
	certFile string
	keyFile  string
	sync.RWMutex
	// deadlock.RWMutex // for deadlock testing/issues
}

var instance *Config

func GetInstance() *Config {
	if instance == nil {
		panic("config not initialized")
	}
	return instance
}

func (i *Config) load(f string) error {
	if err := i.main.Load(file.Provider(f), yaml.Parser()); err != nil {
		return err
	}

	i.filePath = f
	return nil
}

func (i *Config) IsNewSystem() bool {
	return i.isNewSystem
}

func (i *Config) SetConfigFile(fn string) {
	i.Lock()
	defer i.Unlock()
	i.filePath = fn
}

func (i *Config) InitTLS() {
	configDirectory := i.GetConfigPath()
	tlsPaths := []string{
		configDirectory,
		paths.GetStashHomeDirectory(),
	}

	i.certFile = i.getString(sslCertPath)
	if i.certFile == "" {
		// Look for default file
		i.certFile = fsutil.FindInPaths(tlsPaths, "stash.crt")
	}

	i.keyFile = i.getString(sslKeyPath)
	if i.keyFile == "" {
		// Look for default file
		i.keyFile = fsutil.FindInPaths(tlsPaths, "stash.key")
	}
}

func (i *Config) GetTLSFiles() (certFile, keyFile string) {
	return i.certFile, i.keyFile
}

func (i *Config) HasTLSConfig() bool {
	certFile, keyFile := i.GetTLSFiles()
	return certFile != "" && keyFile != ""
}

func (i *Config) GetNoBrowser() bool {
	return i.getBool(NoBrowser)
}

func (i *Config) GetNotificationsEnabled() bool {
	return i.getBool(NotificationsEnabled)
}

// GetShowOneTimeMovedNotification shows whether a small notification to inform the user that Stash
// will no longer show a terminal window, and instead will be available in the tray, should be shown.
// It is true when an existing system is started after upgrading, and set to false forever after it is shown.
func (i *Config) GetShowOneTimeMovedNotification() bool {
	return i.getBool(ShowOneTimeMovedNotification)
}

// these methods are intended to ensure type safety (ie no primitive pointers)
func (i *Config) SetBool(key string, value bool) {
	i.SetInterface(key, value)
}

func (i *Config) SetString(key string, value string) {
	i.SetInterface(key, value)
}

func (i *Config) SetInt(key string, value int) {
	i.SetInterface(key, value)
}

func (i *Config) SetFloat(key string, value float64) {
	i.SetInterface(key, value)
}

func (i *Config) SetInterface(key string, value interface{}) {
	i.Lock()
	defer i.Unlock()

	i.set(key, value)
}

func (i *Config) set(key string, value interface{}) {
	// assumes lock held

	// default behaviour for Set is to merge the value
	// we want to replace it
	i.main.Delete(key)

	if value == nil {
		return
	}

	// test for nil interface as well
	refVal := reflect.ValueOf(value)
	if refVal.Kind() == reflect.Ptr && refVal.IsNil() {
		return
	}

	_ = i.main.Set(key, value)
}

func (i *Config) SetDefault(key string, value interface{}) {
	i.Lock()
	defer i.Unlock()

	i.setDefault(key, value)
}

func (i *Config) setDefault(key string, value interface{}) {
	if !i.main.Exists(key) {
		i.set(key, value)
	}
}

func (i *Config) SetPassword(value string) {
	// if blank, don't bother hashing; we want it to be blank
	if value == "" {
		i.SetString(Password, "")
	} else {
		i.SetString(Password, hashPassword(value))
	}
}

func (i *Config) Write() error {
	i.Lock()
	defer i.Unlock()

	data, err := i.marshal()
	if err != nil {
		return err
	}

	return os.WriteFile(i.filePath, data, 0640)
}

func (i *Config) Marshal() ([]byte, error) {
	i.RLock()
	defer i.RUnlock()

	return i.marshal()
}

func (i *Config) marshal() ([]byte, error) {
	return i.main.Marshal(yaml.Parser())
}

// FileEnvSet returns true if the configuration file environment parameter
// is set.
func FileEnvSet() bool {
	return os.Getenv("STASH_CONFIG_FILE") != ""
}

// GetConfigFile returns the full path to the used configuration file.
func (i *Config) GetConfigFile() string {
	i.RLock()
	defer i.RUnlock()
	return i.filePath
}

// GetConfigPath returns the path of the directory containing the used
// configuration file.
func (i *Config) GetConfigPath() string {
	return filepath.Dir(i.GetConfigFile())
}

// GetConfigPathAbs returns the path of the directory containing the used
// configuration file, resolved to an absolute path. Returns the return value
// of GetConfigPath if the path cannot be made into an absolute path.
func (i *Config) GetConfigPathAbs() string {
	p := filepath.Dir(i.GetConfigFile())

	ret, _ := filepath.Abs(p)
	if ret == "" {
		return p
	}

	return ret
}

// GetDefaultDatabaseFilePath returns the default database filename,
// which is located in the same directory as the config file.
func (i *Config) GetDefaultDatabaseFilePath() string {
	return filepath.Join(i.GetConfigPath(), "stash-go.sqlite")
}

// forKey returns the Koanf instance that should be used to get the provided
// key. Returns the overrides instance if the key exists there, otherwise it
// returns the main instance. Assumes read lock held.
func (i *Config) forKey(key string) *koanf.Koanf {
	v := i.main
	if i.overrides.Exists(key) {
		v = i.overrides
	}

	return v
}

// viper returns the viper instance that has the key set. Returns nil
// if no instance has the key. Assumes read lock held.
func (i *Config) with(key string) *koanf.Koanf {
	v := i.forKey(key)

	if v.Exists(key) {
		return v
	}

	return nil
}

func (i *Config) HasOverride(key string) bool {
	i.RLock()
	defer i.RUnlock()

	return i.overrides.Exists(key)
}

// These functions wrap the equivalent viper functions, checking the override
// instance first, then the main instance.

func (i *Config) unmarshalKey(key string, rawVal interface{}) error {
	i.RLock()
	defer i.RUnlock()

	return i.forKey(key).Unmarshal(key, rawVal)
}

func (i *Config) getStringSlice(key string) []string {
	i.RLock()
	defer i.RUnlock()

	return i.forKey(key).Strings(key)
}

func (i *Config) getString(key string) string {
	i.RLock()
	defer i.RUnlock()

	return i.forKey(key).String(key)
}

func (i *Config) getBool(key string) bool {
	i.RLock()
	defer i.RUnlock()

	return i.forKey(key).Bool(key)
}

func (i *Config) getBoolDefault(key string, def bool) bool {
	i.RLock()
	defer i.RUnlock()

	ret := def
	v := i.forKey(key)
	if v.Exists(key) {
		ret = v.Bool(key)
	}
	return ret
}

func (i *Config) getInt(key string) int {
	i.RLock()
	defer i.RUnlock()

	return i.forKey(key).Int(key)
}

func (i *Config) getFloat64(key string) float64 {
	i.RLock()
	defer i.RUnlock()

	return i.forKey(key).Float64(key)
}

func (i *Config) getStringMapString(key string) map[string]string {
	i.RLock()
	defer i.RUnlock()

	ret := i.forKey(key).StringMap(key)

	// GetStringMapString returns an empty map regardless of whether the
	// key exists or not.
	if len(ret) == 0 {
		return nil
	}

	return ret
}

// GetStathPaths returns the configured stash library paths.
// Works opposite to the usual case - it will return the override
// value only if the main value is not set.
func (i *Config) GetStashPaths() StashConfigs {
	i.RLock()
	defer i.RUnlock()

	var ret StashConfigs

	v := i.main
	if !v.Exists(Stash) {
		v = i.overrides
	}

	if err := v.Unmarshal(Stash, &ret); err != nil || len(ret) == 0 {
		// fallback to legacy format
		ss := v.Strings(Stash)
		ret = nil
		for _, path := range ss {
			toAdd := &StashConfig{
				Path: path,
			}
			ret = append(ret, toAdd)
		}
	}

	return ret
}

func (i *Config) GetCachePath() string {
	return i.getString(Cache)
}

func (i *Config) GetGeneratedPath() string {
	return i.getString(Generated)
}

func (i *Config) GetBlobsPath() string {
	return i.getString(BlobsPath)
}

// GetExtraBlobsPaths returns extra blobs paths.
// For developer/advanced use only.
func (i *Config) GetExtraBlobsPaths() []string {
	return i.getStringSlice(ExtraBlobsPaths)
}

func (i *Config) GetBlobsStorage() BlobsStorageType {
	ret := BlobsStorageType(i.getString(BlobsStorage))

	if !ret.IsValid() {
		// default to database storage
		// for legacy systems this is probably the safer option
		ret = BlobStorageTypeDatabase
	}

	return ret
}

func (i *Config) GetMetadataPath() string {
	return i.getString(Metadata)
}

func (i *Config) GetDatabasePath() string {
	return i.getString(Database)
}

func (i *Config) GetBackupDirectoryPath() string {
	return i.getString(BackupDirectoryPath)
}

func (i *Config) GetBackupDirectoryPathOrDefault() string {
	ret := i.GetBackupDirectoryPath()
	if ret == "" {
		return i.GetConfigPath()
	}

	return ret
}

// GetFFMpegPath returns the path to the FFMpeg executable.
// If empty, stash will attempt to resolve it from the path.
func (i *Config) GetFFMpegPath() string {
	return i.getString(FFMpegPath)
}

// GetFFProbePath returns the path to the FFProbe executable.
// If empty, stash will attempt to resolve it from the path.
func (i *Config) GetFFProbePath() string {
	return i.getString(FFProbePath)
}

func (i *Config) GetJWTSignKey() []byte {
	return []byte(i.getString(JWTSignKey))
}

func (i *Config) GetSessionStoreKey() []byte {
	return []byte(i.getString(SessionStoreKey))
}

func (i *Config) GetDefaultScrapersPath() string {
	// default to the same directory as the config file
	fn := filepath.Join(i.GetConfigPath(), "scrapers")

	return fn
}

func (i *Config) GetExcludes() []string {
	return i.getStringSlice(Exclude)
}

func (i *Config) GetImageExcludes() []string {
	return i.getStringSlice(ImageExclude)
}

func (i *Config) GetVideoExtensions() []string {
	ret := i.getStringSlice(VideoExtensions)
	if len(ret) == 0 {
		ret = defaultVideoExtensions
	}
	return ret
}

func (i *Config) GetImageExtensions() []string {
	ret := i.getStringSlice(ImageExtensions)
	if len(ret) == 0 {
		ret = defaultImageExtensions
	}
	return ret
}

func (i *Config) GetGalleryExtensions() []string {
	ret := i.getStringSlice(GalleryExtensions)
	if len(ret) == 0 {
		ret = defaultGalleryExtensions
	}
	return ret
}

func (i *Config) GetCreateGalleriesFromFolders() bool {
	return i.getBool(CreateGalleriesFromFolders)
}

func (i *Config) GetLanguage() string {
	ret := i.getString(Language)

	// default to English
	if ret == "" {
		return "en-US"
	}

	return ret
}

// IsCalculateMD5 returns true if MD5 checksums should be generated for
// scene video files.
func (i *Config) IsCalculateMD5() bool {
	return i.getBool(CalculateMD5)
}

// GetVideoFileNamingAlgorithm returns what hash algorithm should be used for
// naming generated scene video files.
func (i *Config) GetVideoFileNamingAlgorithm() models.HashAlgorithm {
	ret := i.getString(VideoFileNamingAlgorithm)

	// default to oshash
	if ret == "" {
		return models.HashAlgorithmOshash
	}

	return models.HashAlgorithm(ret)
}

func (i *Config) GetSequentialScanning() bool {
	return i.getBool(SequentialScanning)
}

func (i *Config) GetGalleryCoverRegex() string {
	var regexString = i.getString(GalleryCoverRegex)

	_, err := regexp.Compile(regexString)
	if err != nil {
		logger.Warnf("Gallery cover regex '%v' invalid, reverting to default.", regexString)
		return galleryCoverRegexDefault
	}

	return regexString
}

func (i *Config) GetScrapersPath() string {
	return i.getString(ScrapersPath)
}

func (i *Config) GetScraperUserAgent() string {
	return i.getString(ScraperUserAgent)
}

// GetScraperCDPPath gets the path to the Chrome executable or remote address
// to an instance of Chrome.
func (i *Config) GetScraperCDPPath() string {
	return i.getString(ScraperCDPPath)
}

// GetScraperCertCheck returns true if the scraper should check for insecure
// certificates when fetching an image or a page.
func (i *Config) GetScraperCertCheck() bool {
	return i.getBoolDefault(ScraperCertCheck, true)
}

func (i *Config) GetScraperExcludeTagPatterns() []string {
	return i.getStringSlice(ScraperExcludeTagPatterns)
}

func (i *Config) GetStashBoxes() []*models.StashBox {
	var boxes []*models.StashBox
	if err := i.unmarshalKey(StashBoxes, &boxes); err != nil {
		logger.Warnf("error in unmarshalkey: %v", err)
	}

	return boxes
}

func (i *Config) GetDefaultPluginsPath() string {
	// default to the same directory as the config file
	fn := filepath.Join(i.GetConfigPath(), "plugins")

	return fn
}

func (i *Config) GetPluginsPath() string {
	return i.getString(PluginsPath)
}

func (i *Config) GetAllPluginConfiguration() map[string]map[string]interface{} {
	i.RLock()
	defer i.RUnlock()

	ret := make(map[string]map[string]interface{})

	v := i.forKey(PluginsSetting)

	sub := v.Cut(PluginsSetting)
	if sub == nil {
		return ret
	}

	for plugin := range sub.Raw() {
		ret[plugin] = sub.Cut(plugin).Raw()
	}

	return ret
}

func (i *Config) GetPluginConfiguration(pluginID string) map[string]interface{} {
	i.RLock()
	defer i.RUnlock()

	key := PluginsSettingPrefix + pluginID

	return i.forKey(key).Cut(key).Raw()
}

// SetPluginConfiguration sets the configuration for a plugin.
// It will overwrite any existing configuration.
func (i *Config) SetPluginConfiguration(pluginID string, v map[string]interface{}) {
	i.Lock()
	defer i.Unlock()

	key := PluginsSettingPrefix + pluginID

	i.set(key, v)
}

func (i *Config) GetDisabledPlugins() []string {
	return i.getStringSlice(DisabledPlugins)
}

func (i *Config) GetPythonPath() string {
	return i.getString(PythonPath)
}

func (i *Config) GetHost() string {
	ret := i.getString(Host)
	if ret == "" {
		ret = hostDefault
	}

	return ret
}

func (i *Config) GetPort() int {
	ret := i.getInt(Port)
	if ret == 0 {
		ret = portDefault
	}

	return ret
}

func (i *Config) GetThemeColor() string {
	return i.getString(ThemeColor)
}

func (i *Config) GetExternalHost() string {
	return i.getString(ExternalHost)
}

// GetPreviewSegmentDuration returns the duration of a single segment in a
// scene preview file, in seconds.
func (i *Config) GetPreviewSegmentDuration() float64 {
	return i.getFloat64(PreviewSegmentDuration)
}

// GetParallelTasks returns the number of parallel tasks that should be started
// by scan or generate task.
func (i *Config) GetParallelTasks() int {
	return i.getInt(ParallelTasks)
}

func (i *Config) GetParallelTasksWithAutoDetection() int {
	parallelTasks := i.getInt(ParallelTasks)
	if parallelTasks <= 0 {
		parallelTasks = (runtime.NumCPU() / 4) + 1
	}
	return parallelTasks
}

func (i *Config) GetPreviewAudio() bool {
	return i.getBool(PreviewAudio)
}

// GetPreviewSegments returns the amount of segments in a scene preview file.
func (i *Config) GetPreviewSegments() int {
	return i.getInt(PreviewSegments)
}

// GetPreviewExcludeStart returns the configuration setting string for
// excluding the start of scene videos for preview generation. This can
// be in two possible formats. A float value is interpreted as the amount
// of seconds to exclude from the start of the video before it is included
// in the preview. If the value is suffixed with a '%' character (for example
// '2%'), then it is interpreted as a proportion of the total video duration.
func (i *Config) GetPreviewExcludeStart() string {
	return i.getString(PreviewExcludeStart)
}

// GetPreviewExcludeEnd returns the configuration setting string for
// excluding the end of scene videos for preview generation. A float value
// is interpreted as the amount of seconds to exclude from the end of the video
// when generating previews. If the value is suffixed with a '%' character,
// then it is interpreted as a proportion of the total video duration.
func (i *Config) GetPreviewExcludeEnd() string {
	return i.getString(PreviewExcludeEnd)
}

// GetPreviewPreset returns the preset when generating previews. Defaults to
// Slow.
func (i *Config) GetPreviewPreset() models.PreviewPreset {
	ret := i.getString(PreviewPreset)

	// default to slow
	if ret == "" {
		return models.PreviewPresetSlow
	}

	return models.PreviewPreset(ret)
}

func (i *Config) GetTranscodeHardwareAcceleration() bool {
	return i.getBool(TranscodeHardwareAcceleration)
}

func (i *Config) GetMaxTranscodeSize() models.StreamingResolutionEnum {
	ret := i.getString(MaxTranscodeSize)

	// default to original
	if ret == "" {
		return models.StreamingResolutionEnumOriginal
	}

	return models.StreamingResolutionEnum(ret)
}

func (i *Config) GetMaxStreamingTranscodeSize() models.StreamingResolutionEnum {
	ret := i.getString(MaxStreamingTranscodeSize)

	// default to original
	if ret == "" {
		return models.StreamingResolutionEnumOriginal
	}

	return models.StreamingResolutionEnum(ret)
}

func (i *Config) GetTranscodeInputArgs() []string {
	return i.getStringSlice(TranscodeInputArgs)
}

func (i *Config) GetTranscodeOutputArgs() []string {
	return i.getStringSlice(TranscodeOutputArgs)
}

func (i *Config) GetLiveTranscodeInputArgs() []string {
	return i.getStringSlice(LiveTranscodeInputArgs)
}

func (i *Config) GetLiveTranscodeOutputArgs() []string {
	return i.getStringSlice(LiveTranscodeOutputArgs)
}

func (i *Config) GetDrawFunscriptHeatmapRange() bool {
	return i.getBoolDefault(DrawFunscriptHeatmapRange, drawFunscriptHeatmapRangeDefault)
}

// IsWriteImageThumbnails returns true if image thumbnails should be written
// to disk after generating on the fly.
func (i *Config) IsWriteImageThumbnails() bool {
	return i.getBool(WriteImageThumbnails)
}

func (i *Config) IsCreateImageClipsFromVideos() bool {
	return i.getBool(CreateImageClipsFromVideos)
}

func (i *Config) GetAPIKey() string {
	return i.getString(ApiKey)
}

func (i *Config) GetUsername() string {
	return i.getString(Username)
}

func (i *Config) GetPasswordHash() string {
	return i.getString(Password)
}

func (i *Config) GetCredentials() (string, string) {
	if i.HasCredentials() {
		return i.getString(Username), i.getString(Password)
	}

	return "", ""
}

func (i *Config) HasCredentials() bool {
	username := i.getString(Username)
	pwHash := i.getString(Password)

	return username != "" && pwHash != ""
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hash)
}

func (i *Config) ValidateCredentials(username string, password string) bool {
	if !i.HasCredentials() {
		// don't need to authenticate if no credentials saved
		return true
	}

	authUser, authPWHash := i.GetCredentials()

	err := bcrypt.CompareHashAndPassword([]byte(authPWHash), []byte(password))

	return username == authUser && err == nil
}

var stashBoxRe = regexp.MustCompile("^http.*graphql$")

type StashBoxInput struct {
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"api_key"`
	Name     string `json:"name"`
}

func (i *Config) ValidateStashBoxes(boxes []*StashBoxInput) error {
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
func (i *Config) GetMaxSessionAge() int {
	i.RLock()
	defer i.RUnlock()

	ret := DefaultMaxSessionAge
	v := i.forKey(MaxSessionAge)
	if v.Exists(MaxSessionAge) {
		ret = v.Int(MaxSessionAge)
	}

	return ret
}

// GetCustomServedFolders gets the map of custom paths to their applicable
// filesystem locations
func (i *Config) GetCustomServedFolders() utils.URLMap {
	return i.getStringMapString(CustomServedFolders)
}

func (i *Config) GetUILocation() string {
	if ret := i.getString(UILocation); ret != "" {
		return ret
	}

	return i.getString(LegacyCustomUILocation)
}

// Interface options
func (i *Config) GetMenuItems() []string {
	i.RLock()
	defer i.RUnlock()
	v := i.forKey(MenuItems)
	if v.Exists(MenuItems) {
		return v.Strings(MenuItems)
	}
	return defaultMenuItems
}

func (i *Config) GetSoundOnPreview() bool {
	return i.getBool(SoundOnPreview)
}

func (i *Config) GetWallShowTitle() bool {
	i.RLock()
	defer i.RUnlock()

	ret := defaultWallShowTitle
	v := i.forKey(WallShowTitle)
	if v.Exists(WallShowTitle) {
		ret = v.Bool(WallShowTitle)
	}
	return ret
}

func (i *Config) GetCustomPerformerImageLocation() string {
	return i.getString(CustomPerformerImageLocation)
}

func (i *Config) GetWallPlayback() string {
	i.RLock()
	defer i.RUnlock()

	ret := defaultWallPlayback
	v := i.forKey(WallPlayback)
	if v.Exists(WallPlayback) {
		ret = v.String(WallPlayback)
	}

	return ret
}

func (i *Config) GetShowScrubber() bool {
	return i.getBoolDefault(ShowScrubber, showScrubberDefault)
}

func (i *Config) GetMaximumLoopDuration() int {
	return i.getInt(MaximumLoopDuration)
}

func (i *Config) GetAutostartVideo() bool {
	return i.getBool(AutostartVideo)
}

func (i *Config) GetAutostartVideoOnPlaySelected() bool {
	return i.getBoolDefault(AutostartVideoOnPlaySelected, autostartVideoOnPlaySelectedDefault)
}

func (i *Config) GetContinuePlaylistDefault() bool {
	return i.getBool(ContinuePlaylistDefault)
}

func (i *Config) GetShowStudioAsText() bool {
	return i.getBool(ShowStudioAsText)
}

func (i *Config) getSlideshowDelay() int {
	// assume have lock

	ret := defaultImageLightboxSlideshowDelay
	v := i.forKey(ImageLightboxSlideshowDelay)
	if v.Exists(ImageLightboxSlideshowDelay) {
		ret = v.Int(ImageLightboxSlideshowDelay)
	} else {
		// fallback to old location
		v := i.forKey(legacyImageLightboxSlideshowDelay)
		if v.Exists(legacyImageLightboxSlideshowDelay) {
			ret = v.Int(legacyImageLightboxSlideshowDelay)
		}
	}

	return ret
}

func (i *Config) GetImageLightboxOptions() ConfigImageLightboxResult {
	i.RLock()
	defer i.RUnlock()

	delay := i.getSlideshowDelay()

	ret := ConfigImageLightboxResult{
		SlideshowDelay: &delay,
	}

	if v := i.with(ImageLightboxDisplayModeKey); v != nil {
		mode := ImageLightboxDisplayMode(v.String(ImageLightboxDisplayModeKey))
		ret.DisplayMode = &mode
	}
	if v := i.with(ImageLightboxScaleUp); v != nil {
		value := v.Bool(ImageLightboxScaleUp)
		ret.ScaleUp = &value
	}
	if v := i.with(ImageLightboxResetZoomOnNav); v != nil {
		value := v.Bool(ImageLightboxResetZoomOnNav)
		ret.ResetZoomOnNav = &value
	}
	if v := i.with(ImageLightboxScrollModeKey); v != nil {
		mode := ImageLightboxScrollMode(v.String(ImageLightboxScrollModeKey))
		ret.ScrollMode = &mode
	}
	if v := i.with(ImageLightboxScrollAttemptsBeforeChange); v != nil {
		ret.ScrollAttemptsBeforeChange = v.Int(ImageLightboxScrollAttemptsBeforeChange)
	}

	return ret
}

func (i *Config) GetDisableDropdownCreate() *ConfigDisableDropdownCreate {
	return &ConfigDisableDropdownCreate{
		Performer: i.getBool(DisableDropdownCreatePerformer),
		Studio:    i.getBool(DisableDropdownCreateStudio),
		Tag:       i.getBool(DisableDropdownCreateTag),
		Movie:     i.getBool(DisableDropdownCreateMovie),
	}
}

func (i *Config) GetUIConfiguration() map[string]interface{} {
	i.RLock()
	defer i.RUnlock()

	return i.forKey(UI).Cut(UI).Raw()
}

func (i *Config) SetUIConfiguration(v map[string]interface{}) {
	i.Lock()
	defer i.Unlock()

	i.set(UI, v)
}

func (i *Config) GetCSSPath() string {
	// use custom.css in the same directory as the config file
	configFileUsed := i.GetConfigFile()
	configDir := filepath.Dir(configFileUsed)

	fn := filepath.Join(configDir, "custom.css")

	return fn
}

func (i *Config) GetCSS() string {
	fn := i.GetCSSPath()

	exists, _ := fsutil.FileExists(fn)
	if !exists {
		return ""
	}

	buf, err := os.ReadFile(fn)

	if err != nil {
		return ""
	}

	return string(buf)
}

func (i *Config) SetCSS(css string) {
	fn := i.GetCSSPath()
	i.Lock()
	defer i.Unlock()

	buf := []byte(css)

	if err := os.WriteFile(fn, buf, 0777); err != nil {
		logger.Warnf("error while writing %v bytes to %v: %v", len(buf), fn, err)
	}
}

func (i *Config) GetCSSEnabled() bool {
	return i.getBool(CSSEnabled)
}

func (i *Config) GetJavascriptPath() string {
	// use custom.js in the same directory as the config file
	configFileUsed := i.GetConfigFile()
	configDir := filepath.Dir(configFileUsed)

	fn := filepath.Join(configDir, "custom.js")

	return fn
}

func (i *Config) GetJavascript() string {
	fn := i.GetJavascriptPath()

	exists, _ := fsutil.FileExists(fn)
	if !exists {
		return ""
	}

	buf, err := os.ReadFile(fn)

	if err != nil {
		return ""
	}

	return string(buf)
}

func (i *Config) SetJavascript(javascript string) {
	fn := i.GetJavascriptPath()
	i.Lock()
	defer i.Unlock()

	buf := []byte(javascript)

	if err := os.WriteFile(fn, buf, 0777); err != nil {
		logger.Warnf("error while writing %v bytes to %v: %v", len(buf), fn, err)
	}
}

func (i *Config) GetJavascriptEnabled() bool {
	return i.getBool(JavascriptEnabled)
}

func (i *Config) GetCustomLocalesPath() string {
	// use custom-locales.json in the same directory as the config file
	configFileUsed := i.GetConfigFile()
	configDir := filepath.Dir(configFileUsed)

	fn := filepath.Join(configDir, "custom-locales.json")

	return fn
}

func (i *Config) GetCustomLocales() string {
	fn := i.GetCustomLocalesPath()

	exists, _ := fsutil.FileExists(fn)
	if !exists {
		return ""
	}

	buf, err := os.ReadFile(fn)

	if err != nil {
		return ""
	}

	return string(buf)
}

func (i *Config) SetCustomLocales(customLocales string) {
	fn := i.GetCustomLocalesPath()
	i.Lock()
	defer i.Unlock()

	buf := []byte(customLocales)

	if err := os.WriteFile(fn, buf, 0777); err != nil {
		logger.Warnf("error while writing %v bytes to %v: %v", len(buf), fn, err)
	}
}

func (i *Config) GetCustomLocalesEnabled() bool {
	return i.getBool(CustomLocalesEnabled)
}

func (i *Config) GetHandyKey() string {
	return i.getString(HandyKey)
}

func (i *Config) GetFunscriptOffset() int {
	return i.getInt(FunscriptOffset)
}

func (i *Config) GetUseStashHostedFunscript() bool {
	return i.getBoolDefault(UseStashHostedFunscript, useStashHostedFunscriptDefault)
}

func (i *Config) GetDeleteFileDefault() bool {
	return i.getBool(DeleteFileDefault)
}

func (i *Config) GetDeleteGeneratedDefault() bool {
	return i.getBoolDefault(DeleteGeneratedDefault, deleteGeneratedDefaultDefault)
}

// GetDefaultIdentifySettings returns the default Identify task settings.
// Returns nil if the settings could not be unmarshalled, or if it
// has not been set.
func (i *Config) GetDefaultIdentifySettings() *identify.Options {
	i.RLock()
	defer i.RUnlock()
	v := i.forKey(DefaultIdentifySettings)

	if v.Exists(DefaultIdentifySettings) && v.Get(DefaultIdentifySettings) != nil {
		var ret identify.Options

		if err := v.Unmarshal(DefaultIdentifySettings, &ret); err != nil {
			return nil
		}
		return &ret
	}

	return nil
}

// GetDefaultScanSettings returns the default Scan task settings.
// Returns nil if the settings could not be unmarshalled, or if it
// has not been set.
func (i *Config) GetDefaultScanSettings() *ScanMetadataOptions {
	i.RLock()
	defer i.RUnlock()
	v := i.forKey(DefaultScanSettings)

	if v.Exists(DefaultScanSettings) && v.Get(DefaultScanSettings) != nil {
		var ret ScanMetadataOptions
		if err := v.Unmarshal(DefaultScanSettings, &ret); err != nil {
			return nil
		}
		return &ret
	}

	return nil
}

// GetDefaultAutoTagSettings returns the default Scan task settings.
// Returns nil if the settings could not be unmarshalled, or if it
// has not been set.
func (i *Config) GetDefaultAutoTagSettings() *AutoTagMetadataOptions {
	i.RLock()
	defer i.RUnlock()
	v := i.forKey(DefaultAutoTagSettings)

	if v.Exists(DefaultAutoTagSettings) {
		var ret AutoTagMetadataOptions
		if err := v.Unmarshal(DefaultAutoTagSettings, &ret); err != nil {
			return nil
		}
		return &ret
	}

	return nil
}

// GetDefaultGenerateSettings returns the default Scan task settings.
// Returns nil if the settings could not be unmarshalled, or if it
// has not been set.
func (i *Config) GetDefaultGenerateSettings() *models.GenerateMetadataOptions {
	i.RLock()
	defer i.RUnlock()
	v := i.forKey(DefaultGenerateSettings)

	if v.Exists(DefaultGenerateSettings) {
		var ret models.GenerateMetadataOptions
		if err := v.Unmarshal(DefaultGenerateSettings, &ret); err != nil {
			return nil
		}
		return &ret
	}

	return nil
}

// GetDangerousAllowPublicWithoutAuth determines if the security feature is enabled.
// See https://docs.stashapp.cc/networking/authentication-required-when-accessing-stash-from-the-internet
func (i *Config) GetDangerousAllowPublicWithoutAuth() bool {
	return i.getBool(dangerousAllowPublicWithoutAuth)
}

// GetSecurityTripwireAccessedFromPublicInternet returns a public IP address if stash
// has been accessed from the public internet, with no auth enabled, and
// DangerousAllowPublicWithoutAuth disabled. Returns an empty string otherwise.
func (i *Config) GetSecurityTripwireAccessedFromPublicInternet() string {
	return i.getString(SecurityTripwireAccessedFromPublicInternet)
}

// GetDLNAServerName returns the visible name of the DLNA server. If empty,
// "stash" will be used.
func (i *Config) GetDLNAServerName() string {
	return i.getString(DLNAServerName)
}

// GetDLNADefaultEnabled returns true if the DLNA is enabled by default.
func (i *Config) GetDLNADefaultEnabled() bool {
	return i.getBool(DLNADefaultEnabled)
}

// GetDLNADefaultIPWhitelist returns a list of IP addresses/wildcards that
// are allowed to use the DLNA service.
func (i *Config) GetDLNADefaultIPWhitelist() []string {
	return i.getStringSlice(DLNADefaultIPWhitelist)
}

// GetDLNAInterfaces returns a list of interface names to expose DLNA on. If
// empty, runs on all interfaces.
func (i *Config) GetDLNAInterfaces() []string {
	return i.getStringSlice(DLNAInterfaces)
}

// GetDLNAPort returns the port to run the DLNA server on. If empty, 1338
// will be used.
func (i *Config) GetDLNAPort() int {
	ret := i.getInt(DLNAPort)
	if ret == 0 {
		ret = DLNAPortDefault
	}
	return ret
}

// GetDLNAPortAsString returns the port to run the DLNA server on as a string.
func (i *Config) GetDLNAPortAsString() string {
	return ":" + strconv.Itoa(i.GetDLNAPort())
}

// GetVideoSortOrder returns the sort order to display videos. If
// empty, videos will be sorted by titles.
func (i *Config) GetVideoSortOrder() string {
	ret := i.getString(DLNAVideoSortOrder)
	if ret == "" {
		ret = dlnaVideoSortOrderDefault
	}

	return ret
}

// GetLogFile returns the filename of the file to output logs to.
// An empty string means that file logging will be disabled.
func (i *Config) GetLogFile() string {
	return i.getString(LogFile)
}

// GetLogOut returns true if logging should be output to the terminal
// in addition to writing to a log file. Logging will be output to the
// terminal if file logging is disabled. Defaults to true.
func (i *Config) GetLogOut() bool {
	return i.getBoolDefault(LogOut, defaultLogOut)
}

// GetLogLevel returns the lowest log level to write to the log.
// Should be one of "Debug", "Info", "Warning", "Error"
func (i *Config) GetLogLevel() string {
	value := i.getString(LogLevel)
	if value != "Debug" && value != "Info" && value != "Warning" && value != "Error" && value != "Trace" {
		value = defaultLogLevel
	}

	return value
}

// GetLogAccess returns true if http requests should be logged to the terminal.
// HTTP requests are not logged to the log file. Defaults to true.
func (i *Config) GetLogAccess() bool {
	return i.getBoolDefault(LogAccess, defaultLogAccess)
}

// Max allowed graphql upload size in megabytes
func (i *Config) GetMaxUploadSize() int64 {
	i.RLock()
	defer i.RUnlock()
	ret := int64(1024)

	v := i.forKey(MaxUploadSize)
	if v.Exists(MaxUploadSize) {
		ret = v.Int64(MaxUploadSize)
	}
	return ret << 20
}

// GetProxy returns the url of a http proxy to be used for all outgoing http calls.
func (i *Config) GetProxy() string {
	// Validate format
	reg := regexp.MustCompile(`^((?:socks5h?|https?):\/\/)(([\P{Cc}]+):([\P{Cc}]+)@)?(([a-zA-Z0-9][a-zA-Z0-9.-]*)(:[0-9]{1,5})?)`)
	proxy := i.getString(Proxy)
	if proxy != "" && reg.MatchString(proxy) {
		logger.Debug("Proxy is valid, using it")
		return proxy
	} else if proxy != "" {
		logger.Error("Proxy is invalid, please review your configuration")
		return ""
	}
	return ""
}

// GetProxy returns the url of a http proxy to be used for all outgoing http calls.
func (i *Config) GetNoProxy() string {
	// NoProxy does not require validation, it is validated by the native Go library sufficiently
	return i.getString(NoProxy)
}

// ActivatePublicAccessTripwire sets the security_tripwire_accessed_from_public_internet
// config field to the provided IP address to indicate that stash has been accessed
// from this public IP without authentication.
func (i *Config) ActivatePublicAccessTripwire(requestIP string) error {
	i.SetString(SecurityTripwireAccessedFromPublicInternet, requestIP)
	return i.Write()
}

func (i *Config) getPackageSources(key string) []*models.PackageSource {
	var sources []*models.PackageSource
	if err := i.unmarshalKey(key, &sources); err != nil {
		logger.Warnf("error in unmarshalkey: %v", err)
	}

	return sources
}

func (i *Config) GetPluginPackageSources() []*models.PackageSource {
	return i.getPackageSources(PluginPackageSources)
}

func (i *Config) GetScraperPackageSources() []*models.PackageSource {
	return i.getPackageSources(ScraperPackageSources)
}

type packagePathGetter struct {
	getterFn func() []*models.PackageSource
}

func (g packagePathGetter) GetAllSourcePaths() []string {
	p := g.getterFn()
	var ret []string
	for _, v := range p {
		ret = sliceutil.AppendUnique(ret, v.LocalPath)
	}

	return ret
}

func (g packagePathGetter) GetSourcePath(srcURL string) string {
	p := g.getterFn()

	for _, v := range p {
		if v.URL == srcURL {
			return v.LocalPath
		}
	}

	return ""
}

func (i *Config) GetPluginPackagePathGetter() packagePathGetter {
	return packagePathGetter{
		getterFn: i.GetPluginPackageSources,
	}
}

func (i *Config) GetScraperPackagePathGetter() packagePathGetter {
	return packagePathGetter{
		getterFn: i.GetScraperPackageSources,
	}
}

func (i *Config) Validate() error {
	i.RLock()
	defer i.RUnlock()
	mandatoryPaths := []string{
		Database,
		Generated,
	}

	var missingFields []string

	for _, p := range mandatoryPaths {
		if !i.forKey(p).Exists(p) || i.forKey(p).String(p) == "" {
			missingFields = append(missingFields, p)
		}
	}

	if len(missingFields) > 0 {
		return MissingConfigError{
			missingFields: missingFields,
		}
	}

	if i.GetBlobsStorage() == BlobStorageTypeFilesystem && i.forKey(BlobsPath).String(BlobsPath) == "" {
		return MissingConfigError{
			missingFields: []string{BlobsPath},
		}
	}

	return nil
}

func (i *Config) setDefaultValues() {
	// read data before write lock scope
	defaultDatabaseFilePath := i.GetDefaultDatabaseFilePath()
	defaultScrapersPath := i.GetDefaultScrapersPath()
	defaultPluginsPath := i.GetDefaultPluginsPath()

	i.Lock()
	defer i.Unlock()

	// set the default host and port so that these are written to the config
	// file
	i.setDefault(Host, hostDefault)
	i.setDefault(Port, portDefault)

	i.setDefault(ParallelTasks, parallelTasksDefault)
	i.setDefault(SequentialScanning, SequentialScanningDefault)
	i.setDefault(PreviewSegmentDuration, previewSegmentDurationDefault)
	i.setDefault(PreviewSegments, previewSegmentsDefault)
	i.setDefault(PreviewExcludeStart, previewExcludeStartDefault)
	i.setDefault(PreviewExcludeEnd, previewExcludeEndDefault)
	i.setDefault(PreviewAudio, previewAudioDefault)
	i.setDefault(SoundOnPreview, false)

	i.setDefault(ThemeColor, DefaultThemeColor)

	i.setDefault(WriteImageThumbnails, writeImageThumbnailsDefault)
	i.setDefault(CreateImageClipsFromVideos, createImageClipsFromVideosDefault)

	i.setDefault(Database, defaultDatabaseFilePath)

	i.setDefault(dangerousAllowPublicWithoutAuth, dangerousAllowPublicWithoutAuthDefault)
	i.setDefault(SecurityTripwireAccessedFromPublicInternet, securityTripwireAccessedFromPublicInternetDefault)

	// Set generated to the metadata path for backwards compat
	i.setDefault(Generated, i.main.String(Metadata))

	i.setDefault(NoBrowser, NoBrowserDefault)
	i.setDefault(NotificationsEnabled, NotificationsEnabledDefault)
	i.setDefault(ShowOneTimeMovedNotification, ShowOneTimeMovedNotificationDefault)

	// Set default scrapers and plugins paths
	i.setDefault(ScrapersPath, defaultScrapersPath)
	i.setDefault(PluginsPath, defaultPluginsPath)

	// Set default gallery cover regex
	i.setDefault(GalleryCoverRegex, galleryCoverRegexDefault)

	// Set NoProxy default
	i.setDefault(NoProxy, noProxyDefault)

	// set default package sources
	i.setDefault(PluginPackageSources, []map[string]string{{
		"name":      sourceDefaultName,
		"url":       pluginPackageSourcesDefault,
		"localpath": sourceDefaultPath,
	}})
	i.setDefault(ScraperPackageSources, []map[string]string{{
		"name":      sourceDefaultName,
		"url":       scraperPackageSourcesDefault,
		"localpath": sourceDefaultPath,
	}})
}

// setExistingSystemDefaults sets config options that are new and unset in an existing install,
// but should have a separate default than for brand-new systems, to maintain behavior.
// The config file will not be written.
func (i *Config) setExistingSystemDefaults() {
	i.Lock()
	defer i.Unlock()
	if !i.isNewSystem {
		// Existing systems as of the introduction of auto-browser open should retain existing
		// behavior and not start the browser automatically.
		if !i.main.Exists(NoBrowser) {
			i.set(NoBrowser, true)
		}

		// Existing systems as of the introduction of the taskbar should inform users.
		if !i.main.Exists(ShowOneTimeMovedNotification) {
			i.set(ShowOneTimeMovedNotification, true)
		}
	}
}

// SetInitialConfig fills in missing required config fields. The config file will not be written.
func (i *Config) SetInitialConfig() error {
	// generate some api keys
	const apiKeyLength = 32

	if string(i.GetJWTSignKey()) == "" {
		signKey, err := hash.GenerateRandomKey(apiKeyLength)
		if err != nil {
			return fmt.Errorf("error generating JWTSignKey: %w", err)
		}
		i.SetString(JWTSignKey, signKey)
	}

	if string(i.GetSessionStoreKey()) == "" {
		sessionStoreKey, err := hash.GenerateRandomKey(apiKeyLength)
		if err != nil {
			return fmt.Errorf("error generating session store key: %w", err)
		}
		i.SetString(SessionStoreKey, sessionStoreKey)
	}

	i.setDefaultValues()

	return nil
}

func (i *Config) FinalizeSetup() {
	i.isNewSystem = false
	// i.configUpdates <- 0
}
