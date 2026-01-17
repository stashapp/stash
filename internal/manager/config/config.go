package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"sync"
	// "github.com/sasha-s/go-deadlock" // if you have deadlock issues

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

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

	MaxSessionAge = "max_session_age"

	// SFWContentMode mode config key
	SFWContentMode = "sfw_content_mode"

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

	UseCustomSpriteInterval        = "use_custom_sprite_interval"
	UseCustomSpriteIntervalDefault = false

	SpriteInterval        = "sprite_interval"
	SpriteIntervalDefault = 30

	MinimumSprites        = "minimum_sprites"
	MinimumSpritesDefault = 10

	MaximumSprites        = "maximum_sprites"
	MaximumSpritesDefault = 500

	SpriteScreenshotSize        = "sprite_screenshot_width"
	spriteScreenshotSizeDefault = 160

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
	DisableCustomizations               = "disable_customizations"

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
	ImageLightboxDisableAnimation           = "image_lightbox.disable_animation"

	UI = "ui"

	defaultImageLightboxSlideshowDelay = 5

	DisableDropdownCreatePerformer = "disable_dropdown_create.performer"
	DisableDropdownCreateStudio    = "disable_dropdown_create.studio"
	DisableDropdownCreateTag       = "disable_dropdown_create.tag"
	DisableDropdownCreateMovie     = "disable_dropdown_create.movie"
	DisableDropdownCreateGallery   = "disable_dropdown_create.gallery"

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
	LogFile               = "logfile"
	LogOut                = "logout"
	defaultLogOut         = true
	LogLevel              = "loglevel"
	defaultLogLevel       = "Info"
	LogAccess             = "logaccess"
	defaultLogAccess      = true
	LogFileMaxSize        = "logfile_max_size"
	defaultLogFileMaxSize = 0 // megabytes, default disabled

	// Default settings
	DefaultScanSettings     = "defaults.scan_task"
	DefaultIdentifySettings = "defaults.identify_task"
	DefaultAutoTagSettings  = "defaults.auto_tag_task"
	DefaultGenerateSettings = "defaults.generate_task"

	DeleteFileDefault             = "defaults.delete_file"
	DeleteGeneratedDefault        = "defaults.delete_generated"
	deleteGeneratedDefaultDefault = true

	// Trash/Recycle Bin options
	DeleteTrashPath = "delete_trash_path"

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
	defaultVideoExtensions   = []string{"m4v", "mp4", "mov", "wmv", "avi", "mpg", "mpeg", "rmvb", "rm", "flv", "asf", "mkv", "webm", "f4v"}
	defaultImageExtensions   = []string{"png", "jpg", "jpeg", "gif", "webp", "avif"}
	defaultGalleryExtensions = []string{"zip", "cbz"}
	defaultMenuItems         = []string{"scenes", "images", "groups", "markers", "galleries", "performers", "studios", "tags"}
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

	UserStore *UserStore

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

func (s *Config) load(f string) error {
	if err := s.main.Load(file.Provider(f), yaml.Parser()); err != nil {
		return err
	}

	s.filePath = f
	return nil
}

func (s *Config) IsNewSystem() bool {
	return s.isNewSystem
}

func (s *Config) SetConfigFile(fn string) {
	s.Lock()
	defer s.Unlock()
	s.filePath = fn
}

func (s *Config) InitTLS() {
	configDirectory := s.GetConfigPath()
	tlsPaths := []string{
		configDirectory,
		paths.GetStashHomeDirectory(),
	}

	s.certFile = s.getString(sslCertPath)
	if s.certFile == "" {
		// Look for default file
		s.certFile = fsutil.FindInPaths(tlsPaths, "stash.crt")
	}

	s.keyFile = s.getString(sslKeyPath)
	if s.keyFile == "" {
		// Look for default file
		s.keyFile = fsutil.FindInPaths(tlsPaths, "stash.key")
	}
}

func (s *Config) GetTLSFiles() (certFile, keyFile string) {
	return s.certFile, s.keyFile
}

func (s *Config) HasTLSConfig() bool {
	certFile, keyFile := s.GetTLSFiles()
	return certFile != "" && keyFile != ""
}

func (s *Config) GetNoBrowser() bool {
	return s.getBool(NoBrowser)
}

func (s *Config) GetNotificationsEnabled() bool {
	return s.getBool(NotificationsEnabled)
}

// GetShowOneTimeMovedNotification shows whether a small notification to inform the user that Stash
// will no longer show a terminal window, and instead will be available in the tray, should be shown.
// It is true when an existing system is started after upgrading, and set to false forever after it is shown.
func (s *Config) GetShowOneTimeMovedNotification() bool {
	return s.getBool(ShowOneTimeMovedNotification)
}

// these methods are intended to ensure type safety (ie no primitive pointers)
func (s *Config) SetBool(key string, value bool) {
	s.SetInterface(key, value)
}

func (s *Config) SetString(key string, value string) {
	s.SetInterface(key, value)
}

func (s *Config) SetInt(key string, value int) {
	s.SetInterface(key, value)
}

func (s *Config) SetFloat(key string, value float64) {
	s.SetInterface(key, value)
}

func (s *Config) SetInterface(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

	s.setInterfaceNoLock(key, value)
}
func (s *Config) setInterfaceNoLock(key string, value interface{}) {
	s.set(key, value)
}

func (s *Config) set(key string, value interface{}) {
	// assumes lock held

	// default behaviour for Set is to merge the value
	// we want to replace it
	s.main.Delete(key)

	if value == nil {
		return
	}

	// test for nil interface as well
	refVal := reflect.ValueOf(value)
	if refVal.Kind() == reflect.Ptr && refVal.IsNil() {
		return
	}

	_ = s.main.Set(key, value)
}

func (s *Config) SetDefault(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

	s.setDefault(key, value)
}

func (s *Config) setDefault(key string, value interface{}) {
	if !s.main.Exists(key) {
		s.set(key, value)
	}
}

func (s *Config) SetPassword(value string) {
	// if blank, don't bother hashing; we want it to be blank
	if value == "" {
		s.SetString(Password, "")
	} else {
		s.SetString(Password, hashPassword(value))
	}
}

func (s *Config) Write() error {
	s.Lock()
	defer s.Unlock()

	return s.writeNoLock()
}

func (s *Config) writeNoLock() error {
	data, err := s.marshal()
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0640)
}

func (s *Config) Marshal() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	return s.marshal()
}

func (s *Config) marshal() ([]byte, error) {
	return s.main.Marshal(yaml.Parser())
}

// FileEnvSet returns true if the configuration file environment parameter
// is set.
func FileEnvSet() bool {
	return os.Getenv("STASH_CONFIG_FILE") != ""
}

// GetConfigFile returns the full path to the used configuration file.
func (s *Config) GetConfigFile() string {
	s.RLock()
	defer s.RUnlock()
	return s.filePath
}

// GetConfigPath returns the path of the directory containing the used
// configuration file.
func (s *Config) GetConfigPath() string {
	return filepath.Dir(s.GetConfigFile())
}

// GetConfigPathAbs returns the path of the directory containing the used
// configuration file, resolved to an absolute path. Returns the return value
// of GetConfigPath if the path cannot be made into an absolute path.
func (s *Config) GetConfigPathAbs() string {
	p := filepath.Dir(s.GetConfigFile())

	ret, _ := filepath.Abs(p)
	if ret == "" {
		return p
	}

	return ret
}

// GetDefaultDatabaseFilePath returns the default database filename,
// which is located in the same directory as the config file.
func (s *Config) GetDefaultDatabaseFilePath() string {
	return filepath.Join(s.GetConfigPath(), "stash-go.sqlite")
}

// forKey returns the Koanf instance that should be used to get the provided
// key. Returns the overrides instance if the key exists there, otherwise it
// returns the main instance. Assumes read lock held.
func (s *Config) forKey(key string) *koanf.Koanf {
	v := s.main
	if s.overrides.Exists(key) {
		v = s.overrides
	}

	return v
}

// viper returns the viper instance that has the key set. Returns nil
// if no instance has the key. Assumes read lock held.
func (s *Config) with(key string) *koanf.Koanf {
	v := s.forKey(key)

	if v.Exists(key) {
		return v
	}

	return nil
}

func (s *Config) HasOverride(key string) bool {
	s.RLock()
	defer s.RUnlock()

	return s.overrides.Exists(key)
}

// These functions wrap the equivalent viper functions, checking the override
// instance first, then the main instance.

func (s *Config) unmarshalKey(key string, rawVal interface{}) error {
	s.RLock()
	defer s.RUnlock()

	return s.forKey(key).Unmarshal(key, rawVal)
}

func (s *Config) getStringSlice(key string) []string {
	s.RLock()
	defer s.RUnlock()

	return s.forKey(key).Strings(key)
}

func (s *Config) getString(key string) string {
	s.RLock()
	defer s.RUnlock()

	return s.forKey(key).String(key)
}

func (s *Config) getBool(key string) bool {
	s.RLock()
	defer s.RUnlock()

	return s.forKey(key).Bool(key)
}

func (s *Config) getBoolDefault(key string, def bool) bool {
	s.RLock()
	defer s.RUnlock()

	ret := def
	v := s.forKey(key)
	if v.Exists(key) {
		ret = v.Bool(key)
	}
	return ret
}

func (s *Config) getInt(key string) int {
	s.RLock()
	defer s.RUnlock()

	return s.forKey(key).Int(key)
}

func (s *Config) getFloat64(key string) float64 {
	s.RLock()
	defer s.RUnlock()

	return s.forKey(key).Float64(key)
}

func (s *Config) getStringMapString(key string) map[string]string {
	s.RLock()
	defer s.RUnlock()

	ret := s.forKey(key).StringMap(key)

	// GetStringMapString returns an empty map regardless of whether the
	// key exists or not.
	if len(ret) == 0 {
		return nil
	}

	return ret
}

// GetSFW returns true if SFW mode is enabled.
// Default performer images are changed to more agnostic images when enabled.
func (s *Config) GetSFWContentMode() bool {
	s.RLock()
	defer s.RUnlock()
	return s.getBool(SFWContentMode)
}

// GetStashPaths returns the configured stash library paths.
// Works opposite to the usual case - it will return the override
// value only if the main value is not set.
func (s *Config) GetStashPaths() StashConfigs {
	s.RLock()
	defer s.RUnlock()

	var ret StashConfigs

	v := s.main
	if !v.Exists(Stash) {
		v = s.overrides
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

func (s *Config) GetCachePath() string {
	return s.getString(Cache)
}

func (s *Config) GetGeneratedPath() string {
	return s.getString(Generated)
}

func (s *Config) GetBlobsPath() string {
	return s.getString(BlobsPath)
}

// GetExtraBlobsPaths returns extra blobs paths.
// For developer/advanced use only.
func (s *Config) GetExtraBlobsPaths() []string {
	return s.getStringSlice(ExtraBlobsPaths)
}

func (s *Config) GetBlobsStorage() BlobsStorageType {
	ret := BlobsStorageType(s.getString(BlobsStorage))

	if !ret.IsValid() {
		// default to database storage
		// for legacy systems this is probably the safer option
		ret = BlobStorageTypeDatabase
	}

	return ret
}

func (s *Config) GetMetadataPath() string {
	return s.getString(Metadata)
}

func (s *Config) GetDatabasePath() string {
	return s.getString(Database)
}

func (s *Config) GetBackupDirectoryPath() string {
	return s.getString(BackupDirectoryPath)
}

func (s *Config) GetBackupDirectoryPathOrDefault() string {
	ret := s.GetBackupDirectoryPath()
	if ret == "" {
		// #4915 - default to the same directory as the database
		return filepath.Dir(s.GetDatabasePath())
	}

	return ret
}

// GetFFMpegPath returns the path to the FFMpeg executable.
// If empty, stash will attempt to resolve it from the path.
func (s *Config) GetFFMpegPath() string {
	return s.getString(FFMpegPath)
}

// GetFFProbePath returns the path to the FFProbe executable.
// If empty, stash will attempt to resolve it from the path.
func (s *Config) GetFFProbePath() string {
	return s.getString(FFProbePath)
}

func (s *Config) GetJWTSignKey() []byte {
	return []byte(s.getString(JWTSignKey))
}

func (s *Config) GetSessionStoreKey() []byte {
	return []byte(s.getString(SessionStoreKey))
}

func (s *Config) GetDefaultScrapersPath() string {
	// default to the same directory as the config file
	fn := filepath.Join(s.GetConfigPath(), "scrapers")

	return fn
}

func (s *Config) GetExcludes() []string {
	return s.getStringSlice(Exclude)
}

func (s *Config) GetImageExcludes() []string {
	return s.getStringSlice(ImageExclude)
}

func (s *Config) GetVideoExtensions() []string {
	ret := s.getStringSlice(VideoExtensions)
	if len(ret) == 0 {
		ret = defaultVideoExtensions
	}
	return ret
}

func (s *Config) GetImageExtensions() []string {
	ret := s.getStringSlice(ImageExtensions)
	if len(ret) == 0 {
		ret = defaultImageExtensions
	}
	return ret
}

func (s *Config) GetGalleryExtensions() []string {
	ret := s.getStringSlice(GalleryExtensions)
	if len(ret) == 0 {
		ret = defaultGalleryExtensions
	}
	return ret
}

func (s *Config) GetCreateGalleriesFromFolders() bool {
	return s.getBool(CreateGalleriesFromFolders)
}

func (s *Config) GetLanguage() string {
	ret := s.getString(Language)

	// default to English
	if ret == "" {
		return "en-US"
	}

	return ret
}

// IsCalculateMD5 returns true if MD5 checksums should be generated for
// scene video files.
func (s *Config) IsCalculateMD5() bool {
	return s.getBool(CalculateMD5)
}

// GetVideoFileNamingAlgorithm returns what hash algorithm should be used for
// naming generated scene video files.
func (s *Config) GetVideoFileNamingAlgorithm() models.HashAlgorithm {
	ret := s.getString(VideoFileNamingAlgorithm)

	// default to oshash
	if ret == "" {
		return models.HashAlgorithmOshash
	}

	return models.HashAlgorithm(ret)
}

func (s *Config) GetSequentialScanning() bool {
	return s.getBool(SequentialScanning)
}

func (s *Config) GetGalleryCoverRegex() string {
	var regexString = s.getString(GalleryCoverRegex)

	_, err := regexp.Compile(regexString)
	if err != nil {
		logger.Warnf("Gallery cover regex '%v' invalid, reverting to default.", regexString)
		return galleryCoverRegexDefault
	}

	return regexString
}

func (s *Config) GetScrapersPath() string {
	return s.getString(ScrapersPath)
}

func (s *Config) GetScraperUserAgent() string {
	return s.getString(ScraperUserAgent)
}

// GetScraperCDPPath gets the path to the Chrome executable or remote address
// to an instance of Chrome.
func (s *Config) GetScraperCDPPath() string {
	return s.getString(ScraperCDPPath)
}

// GetScraperCertCheck returns true if the scraper should check for insecure
// certificates when fetching an image or a page.
func (s *Config) GetScraperCertCheck() bool {
	return s.getBoolDefault(ScraperCertCheck, true)
}

func (s *Config) GetScraperExcludeTagPatterns() []string {
	return s.getStringSlice(ScraperExcludeTagPatterns)
}

func (s *Config) GetStashBoxes() []*models.StashBox {
	var boxes []*models.StashBox
	if err := s.unmarshalKey(StashBoxes, &boxes); err != nil {
		logger.Warnf("error in unmarshalkey: %v", err)
	}

	return boxes
}

func (s *Config) GetDefaultPluginsPath() string {
	// default to the same directory as the config file
	fn := filepath.Join(s.GetConfigPath(), "plugins")

	return fn
}

func (s *Config) GetPluginsPath() string {
	return s.getString(PluginsPath)
}

func (s *Config) GetAllPluginConfiguration() map[string]map[string]interface{} {
	s.RLock()
	defer s.RUnlock()

	ret := make(map[string]map[string]interface{})

	v := s.forKey(PluginsSetting)

	sub := v.Cut(PluginsSetting)
	if sub == nil {
		return ret
	}

	for plugin := range sub.Raw() {
		ret[plugin] = sub.Cut(plugin).Raw()
	}

	return ret
}

func (s *Config) GetPluginConfiguration(pluginID string) map[string]interface{} {
	s.RLock()
	defer s.RUnlock()

	key := PluginsSettingPrefix + pluginID

	return s.forKey(key).Cut(key).Raw()
}

// SetPluginConfiguration sets the configuration for a plugin.
// It will overwrite any existing configuration.
func (s *Config) SetPluginConfiguration(pluginID string, v map[string]interface{}) {
	s.Lock()
	defer s.Unlock()

	key := PluginsSettingPrefix + pluginID

	s.set(key, v)
}

func (s *Config) GetDisabledPlugins() []string {
	return s.getStringSlice(DisabledPlugins)
}

func (s *Config) GetPythonPath() string {
	return s.getString(PythonPath)
}

func (s *Config) GetHost() string {
	ret := s.getString(Host)
	if ret == "" {
		ret = hostDefault
	}

	return ret
}

func (s *Config) GetPort() int {
	ret := s.getInt(Port)
	if ret == 0 {
		ret = portDefault
	}

	return ret
}

func (s *Config) GetThemeColor() string {
	return s.getString(ThemeColor)
}

func (s *Config) GetExternalHost() string {
	return s.getString(ExternalHost)
}

// GetPreviewSegmentDuration returns the duration of a single segment in a
// scene preview file, in seconds.
func (s *Config) GetPreviewSegmentDuration() float64 {
	return s.getFloat64(PreviewSegmentDuration)
}

// GetParallelTasks returns the number of parallel tasks that should be started
// by scan or generate task.
func (s *Config) GetParallelTasks() int {
	return s.getInt(ParallelTasks)
}

func (s *Config) GetParallelTasksWithAutoDetection() int {
	parallelTasks := s.getInt(ParallelTasks)
	if parallelTasks <= 0 {
		parallelTasks = (runtime.NumCPU() / 4) + 1
	}
	return parallelTasks
}

// GetUseCustomSpriteInterval returns true if the sprite minimum, maximum, and interval settings
// should be used instead of the default
func (s *Config) GetUseCustomSpriteInterval() bool {
	value := s.getBool(UseCustomSpriteInterval)
	return value
}

// GetSpriteInterval returns the time (in seconds) to be between each scrubber sprite
// A value of 0 indicates that the sprite interval should be automatically determined
// based on the minimum sprite setting.
func (s *Config) GetSpriteInterval() float64 {
	value := s.getFloat64(SpriteInterval)
	return value
}

// GetMinimumSprites returns the minimum number of sprites that have to be generated
// A value of 0 will be overridden with the default of 10.
func (s *Config) GetMinimumSprites() int {
	value := s.getInt(MinimumSprites)
	if value <= 0 {
		return MinimumSpritesDefault
	}
	return value
}

// GetMaximumSprites returns the maximum number of sprites that can be generated
// A value of 0 indicates no maximum.
func (s *Config) GetMaximumSprites() int {
	value := s.getInt(MaximumSprites)
	return value
}

// GetSpriteScreenshotSize returns the required size of the screenshots to be taken
// during sprite generation in pixels. This will be the width for landscape scenes
// and the height for portrait scenes, with the other dimension being scaled to maintain
// the aspect ratio. If the value is less than or equal to 0, the default will be used.
func (s *Config) GetSpriteScreenshotSize() int {
	value := s.getInt(SpriteScreenshotSize)
	if value <= 0 {
		return spriteScreenshotSizeDefault
	}
	return value
}

func (s *Config) GetPreviewAudio() bool {
	return s.getBool(PreviewAudio)
}

// GetPreviewSegments returns the amount of segments in a scene preview file.
func (s *Config) GetPreviewSegments() int {
	return s.getInt(PreviewSegments)
}

// GetPreviewExcludeStart returns the configuration setting string for
// excluding the start of scene videos for preview generation. This can
// be in two possible formats. A float value is interpreted as the amount
// of seconds to exclude from the start of the video before it is included
// in the preview. If the value is suffixed with a '%' character (for example
// '2%'), then it is interpreted as a proportion of the total video duration.
func (s *Config) GetPreviewExcludeStart() string {
	return s.getString(PreviewExcludeStart)
}

// GetPreviewExcludeEnd returns the configuration setting string for
// excluding the end of scene videos for preview generation. A float value
// is interpreted as the amount of seconds to exclude from the end of the video
// when generating previews. If the value is suffixed with a '%' character,
// then it is interpreted as a proportion of the total video duration.
func (s *Config) GetPreviewExcludeEnd() string {
	return s.getString(PreviewExcludeEnd)
}

// GetPreviewPreset returns the preset when generating previews. Defaults to
// Slow.
func (s *Config) GetPreviewPreset() models.PreviewPreset {
	ret := s.getString(PreviewPreset)

	// default to slow
	if ret == "" {
		return models.PreviewPresetSlow
	}

	return models.PreviewPreset(ret)
}

func (s *Config) GetTranscodeHardwareAcceleration() bool {
	return s.getBool(TranscodeHardwareAcceleration)
}

func (s *Config) GetMaxTranscodeSize() models.StreamingResolutionEnum {
	ret := s.getString(MaxTranscodeSize)

	// default to original
	if ret == "" {
		return models.StreamingResolutionEnumOriginal
	}

	return models.StreamingResolutionEnum(ret)
}

func (s *Config) GetMaxStreamingTranscodeSize() models.StreamingResolutionEnum {
	ret := s.getString(MaxStreamingTranscodeSize)

	// default to original
	if ret == "" {
		return models.StreamingResolutionEnumOriginal
	}

	return models.StreamingResolutionEnum(ret)
}

func (s *Config) GetTranscodeInputArgs() []string {
	return s.getStringSlice(TranscodeInputArgs)
}

func (s *Config) GetTranscodeOutputArgs() []string {
	return s.getStringSlice(TranscodeOutputArgs)
}

func (s *Config) GetLiveTranscodeInputArgs() []string {
	return s.getStringSlice(LiveTranscodeInputArgs)
}

func (s *Config) GetLiveTranscodeOutputArgs() []string {
	return s.getStringSlice(LiveTranscodeOutputArgs)
}

func (s *Config) GetDrawFunscriptHeatmapRange() bool {
	return s.getBoolDefault(DrawFunscriptHeatmapRange, drawFunscriptHeatmapRangeDefault)
}

// IsWriteImageThumbnails returns true if image thumbnails should be written
// to disk after generating on the fly.
func (s *Config) IsWriteImageThumbnails() bool {
	return s.getBool(WriteImageThumbnails)
}

func (s *Config) IsCreateImageClipsFromVideos() bool {
	return s.getBool(CreateImageClipsFromVideos)
}

func (s *Config) GetAPIKey() string {
	return s.getString(ApiKey)
}

func stashBoxValidate(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != "" && strings.HasSuffix(u.Path, "/graphql")
}

type StashBoxInput struct {
	Endpoint             string `json:"endpoint"`
	APIKey               string `json:"api_key"`
	Name                 string `json:"name"`
	MaxRequestsPerMinute int    `json:"max_requests_per_minute"`
}

func (s *Config) ValidateStashBoxes(boxes []*StashBoxInput) error {
	isMulti := len(boxes) > 1

	for _, box := range boxes {
		// Validate each stash-box configuration field, return on error
		if box.APIKey == "" {
			return &StashBoxError{msg: "API Key cannot be blank"}
		}

		if box.Endpoint == "" {
			return &StashBoxError{msg: "endpoint cannot be blank"}
		}

		if !stashBoxValidate(box.Endpoint) {
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
func (s *Config) GetMaxSessionAge() int {
	s.RLock()
	defer s.RUnlock()

	ret := DefaultMaxSessionAge
	v := s.forKey(MaxSessionAge)
	if v.Exists(MaxSessionAge) {
		ret = v.Int(MaxSessionAge)
	}

	return ret
}

// GetCustomServedFolders gets the map of custom paths to their applicable
// filesystem locations
func (s *Config) GetCustomServedFolders() utils.URLMap {
	return s.getStringMapString(CustomServedFolders)
}

func (s *Config) GetUILocation() string {
	if ret := s.getString(UILocation); ret != "" {
		return ret
	}

	return s.getString(LegacyCustomUILocation)
}

// Interface options
func (s *Config) GetMenuItems() []string {
	s.RLock()
	defer s.RUnlock()
	v := s.forKey(MenuItems)
	if v.Exists(MenuItems) {
		return v.Strings(MenuItems)
	}
	return defaultMenuItems
}

func (s *Config) GetSoundOnPreview() bool {
	return s.getBool(SoundOnPreview)
}

func (s *Config) GetWallShowTitle() bool {
	s.RLock()
	defer s.RUnlock()

	ret := defaultWallShowTitle
	v := s.forKey(WallShowTitle)
	if v.Exists(WallShowTitle) {
		ret = v.Bool(WallShowTitle)
	}
	return ret
}

func (s *Config) GetCustomPerformerImageLocation() string {
	return s.getString(CustomPerformerImageLocation)
}

func (s *Config) GetWallPlayback() string {
	s.RLock()
	defer s.RUnlock()

	ret := defaultWallPlayback
	v := s.forKey(WallPlayback)
	if v.Exists(WallPlayback) {
		ret = v.String(WallPlayback)
	}

	return ret
}

func (s *Config) GetShowScrubber() bool {
	return s.getBoolDefault(ShowScrubber, showScrubberDefault)
}

func (s *Config) GetMaximumLoopDuration() int {
	return s.getInt(MaximumLoopDuration)
}

func (s *Config) GetAutostartVideo() bool {
	return s.getBool(AutostartVideo)
}

func (s *Config) GetAutostartVideoOnPlaySelected() bool {
	return s.getBoolDefault(AutostartVideoOnPlaySelected, autostartVideoOnPlaySelectedDefault)
}

func (s *Config) GetContinuePlaylistDefault() bool {
	return s.getBool(ContinuePlaylistDefault)
}

func (s *Config) GetShowStudioAsText() bool {
	return s.getBool(ShowStudioAsText)
}

func (s *Config) getSlideshowDelay() int {
	// assume have lock

	ret := defaultImageLightboxSlideshowDelay
	v := s.forKey(ImageLightboxSlideshowDelay)
	if v.Exists(ImageLightboxSlideshowDelay) {
		ret = v.Int(ImageLightboxSlideshowDelay)
	} else {
		// fallback to old location
		v := s.forKey(legacyImageLightboxSlideshowDelay)
		if v.Exists(legacyImageLightboxSlideshowDelay) {
			ret = v.Int(legacyImageLightboxSlideshowDelay)
		}
	}

	return ret
}

func (s *Config) GetImageLightboxOptions() ConfigImageLightboxResult {
	s.RLock()
	defer s.RUnlock()

	delay := s.getSlideshowDelay()

	ret := ConfigImageLightboxResult{
		SlideshowDelay: &delay,
	}

	if v := s.with(ImageLightboxDisplayModeKey); v != nil {
		mode := ImageLightboxDisplayMode(v.String(ImageLightboxDisplayModeKey))
		ret.DisplayMode = &mode
	}
	if v := s.with(ImageLightboxScaleUp); v != nil {
		value := v.Bool(ImageLightboxScaleUp)
		ret.ScaleUp = &value
	}
	if v := s.with(ImageLightboxResetZoomOnNav); v != nil {
		value := v.Bool(ImageLightboxResetZoomOnNav)
		ret.ResetZoomOnNav = &value
	}
	if v := s.with(ImageLightboxScrollModeKey); v != nil {
		mode := ImageLightboxScrollMode(v.String(ImageLightboxScrollModeKey))
		ret.ScrollMode = &mode
	}
	if v := s.with(ImageLightboxScrollAttemptsBeforeChange); v != nil {
		ret.ScrollAttemptsBeforeChange = v.Int(ImageLightboxScrollAttemptsBeforeChange)
	}
	if v := s.with(ImageLightboxDisableAnimation); v != nil {
		value := v.Bool(ImageLightboxDisableAnimation)
		ret.DisableAnimation = &value
	}

	return ret
}

func (s *Config) GetDisableDropdownCreate() *ConfigDisableDropdownCreate {
	return &ConfigDisableDropdownCreate{
		Performer: s.getBool(DisableDropdownCreatePerformer),
		Studio:    s.getBool(DisableDropdownCreateStudio),
		Tag:       s.getBool(DisableDropdownCreateTag),
		Movie:     s.getBool(DisableDropdownCreateMovie),
		Gallery:   s.getBool(DisableDropdownCreateGallery),
	}
}

func (s *Config) GetUIConfiguration() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()

	return s.forKey(UI).Cut(UI).Raw()
}

// GetMinimumPlayPercent returns the minimum percentage of a video that must be
// watched before incrementing the play count. Returns 0 if not configured.
func (s *Config) GetMinimumPlayPercent() int {
	uiConfig := s.GetUIConfiguration()
	if uiConfig == nil {
		return 0
	}
	if val, ok := uiConfig["minimumPlayPercent"]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case int64:
			return int(v)
		}
	}
	return 0
}

func (s *Config) SetUIConfiguration(v map[string]interface{}) {
	s.Lock()
	defer s.Unlock()

	s.set(UI, v)
}

func (s *Config) GetCSSPath() string {
	// use custom.css in the same directory as the config file
	configFileUsed := s.GetConfigFile()
	configDir := filepath.Dir(configFileUsed)

	fn := filepath.Join(configDir, "custom.css")

	return fn
}

func (s *Config) GetCSS() string {
	fn := s.GetCSSPath()

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

func (s *Config) SetCSS(css string) {
	fn := s.GetCSSPath()
	s.Lock()
	defer s.Unlock()

	buf := []byte(css)

	if err := os.WriteFile(fn, buf, 0777); err != nil {
		logger.Warnf("error while writing %v bytes to %v: %v", len(buf), fn, err)
	}
}

func (s *Config) GetCSSEnabled() bool {
	return s.getBool(CSSEnabled)
}

func (s *Config) GetJavascriptPath() string {
	// use custom.js in the same directory as the config file
	configFileUsed := s.GetConfigFile()
	configDir := filepath.Dir(configFileUsed)

	fn := filepath.Join(configDir, "custom.js")

	return fn
}

func (s *Config) GetJavascript() string {
	fn := s.GetJavascriptPath()

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

func (s *Config) SetJavascript(javascript string) {
	fn := s.GetJavascriptPath()
	s.Lock()
	defer s.Unlock()

	buf := []byte(javascript)

	if err := os.WriteFile(fn, buf, 0777); err != nil {
		logger.Warnf("error while writing %v bytes to %v: %v", len(buf), fn, err)
	}
}

func (s *Config) GetJavascriptEnabled() bool {
	return s.getBool(JavascriptEnabled)
}

func (s *Config) GetCustomLocalesPath() string {
	// use custom-locales.json in the same directory as the config file
	configFileUsed := s.GetConfigFile()
	configDir := filepath.Dir(configFileUsed)

	fn := filepath.Join(configDir, "custom-locales.json")

	return fn
}

func (s *Config) GetCustomLocales() string {
	fn := s.GetCustomLocalesPath()

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

func (s *Config) SetCustomLocales(customLocales string) {
	fn := s.GetCustomLocalesPath()
	s.Lock()
	defer s.Unlock()

	buf := []byte(customLocales)

	if err := os.WriteFile(fn, buf, 0777); err != nil {
		logger.Warnf("error while writing %v bytes to %v: %v", len(buf), fn, err)
	}
}

func (s *Config) GetCustomLocalesEnabled() bool {
	return s.getBool(CustomLocalesEnabled)
}

// GetDisableCustomizations returns true if all customizations (plugins, custom CSS,
// custom JavaScript, and custom locales) should be disabled. This is useful for
// troubleshooting issues without permanently disabling individual customizations.
func (s *Config) GetDisableCustomizations() bool {
	return s.getBool(DisableCustomizations)
}

func (s *Config) GetHandyKey() string {
	return s.getString(HandyKey)
}

func (s *Config) GetFunscriptOffset() int {
	return s.getInt(FunscriptOffset)
}

func (s *Config) GetUseStashHostedFunscript() bool {
	return s.getBoolDefault(UseStashHostedFunscript, useStashHostedFunscriptDefault)
}

func (s *Config) GetDeleteFileDefault() bool {
	return s.getBool(DeleteFileDefault)
}

func (s *Config) GetDeleteGeneratedDefault() bool {
	return s.getBoolDefault(DeleteGeneratedDefault, deleteGeneratedDefaultDefault)
}

func (s *Config) GetDeleteTrashPath() string {
	return s.getString(DeleteTrashPath)
}

func (s *Config) SetDeleteTrashPath(value string) {
	s.SetString(DeleteTrashPath, value)
}

// GetDefaultIdentifySettings returns the default Identify task settings.
// Returns nil if the settings could not be unmarshalled, or if it
// has not been set.
func (s *Config) GetDefaultIdentifySettings() *identify.Options {
	s.RLock()
	defer s.RUnlock()
	v := s.forKey(DefaultIdentifySettings)

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
func (s *Config) GetDefaultScanSettings() *ScanMetadataOptions {
	s.RLock()
	defer s.RUnlock()
	v := s.forKey(DefaultScanSettings)

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
func (s *Config) GetDefaultAutoTagSettings() *AutoTagMetadataOptions {
	s.RLock()
	defer s.RUnlock()
	v := s.forKey(DefaultAutoTagSettings)

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
func (s *Config) GetDefaultGenerateSettings() *models.GenerateMetadataOptions {
	s.RLock()
	defer s.RUnlock()
	v := s.forKey(DefaultGenerateSettings)

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
// See https://discourse.stashapp.cc/t/-/1658
func (s *Config) GetDangerousAllowPublicWithoutAuth() bool {
	return s.getBool(dangerousAllowPublicWithoutAuth)
}

// GetSecurityTripwireAccessedFromPublicInternet returns a public IP address if stash
// has been accessed from the public internet, with no auth enabled, and
// DangerousAllowPublicWithoutAuth disabled. Returns an empty string otherwise.
func (s *Config) GetSecurityTripwireAccessedFromPublicInternet() string {
	return s.getString(SecurityTripwireAccessedFromPublicInternet)
}

// GetDLNAServerName returns the visible name of the DLNA server. If empty,
// "stash" will be used.
func (s *Config) GetDLNAServerName() string {
	return s.getString(DLNAServerName)
}

// GetDLNADefaultEnabled returns true if the DLNA is enabled by default.
func (s *Config) GetDLNADefaultEnabled() bool {
	return s.getBool(DLNADefaultEnabled)
}

// GetDLNADefaultIPWhitelist returns a list of IP addresses/wildcards that
// are allowed to use the DLNA service.
func (s *Config) GetDLNADefaultIPWhitelist() []string {
	return s.getStringSlice(DLNADefaultIPWhitelist)
}

// GetDLNAInterfaces returns a list of interface names to expose DLNA on. If
// empty, runs on all interfaces.
func (s *Config) GetDLNAInterfaces() []string {
	return s.getStringSlice(DLNAInterfaces)
}

// GetDLNAPort returns the port to run the DLNA server on. If empty, 1338
// will be used.
func (s *Config) GetDLNAPort() int {
	ret := s.getInt(DLNAPort)
	if ret == 0 {
		ret = DLNAPortDefault
	}
	return ret
}

// GetDLNAPortAsString returns the port to run the DLNA server on as a string.
func (s *Config) GetDLNAPortAsString() string {
	return ":" + strconv.Itoa(s.GetDLNAPort())
}

// GetDLNAActivityTrackingEnabled returns true if DLNA activity tracking is enabled.
// This uses the same "trackActivity" UI setting that controls frontend play history tracking.
// When enabled, scenes played via DLNA will have their play count and duration tracked.
func (s *Config) GetDLNAActivityTrackingEnabled() bool {
	uiConfig := s.GetUIConfiguration()
	if uiConfig == nil {
		return true // Default to enabled
	}
	if val, ok := uiConfig["trackActivity"]; ok {
		if v, ok := val.(bool); ok {
			return v
		}
	}
	return true // Default to enabled
}

// GetVideoSortOrder returns the sort order to display videos. If
// empty, videos will be sorted by titles.
func (s *Config) GetVideoSortOrder() string {
	ret := s.getString(DLNAVideoSortOrder)
	if ret == "" {
		ret = dlnaVideoSortOrderDefault
	}

	return ret
}

// GetLogFile returns the filename of the file to output logs to.
// An empty string means that file logging will be disabled.
func (s *Config) GetLogFile() string {
	return s.getString(LogFile)
}

// GetLogOut returns true if logging should be output to the terminal
// in addition to writing to a log file. Logging will be output to the
// terminal if file logging is disabled. Defaults to true.
func (s *Config) GetLogOut() bool {
	return s.getBoolDefault(LogOut, defaultLogOut)
}

// GetLogLevel returns the lowest log level to write to the log.
// Should be one of "Debug", "Info", "Warning", "Error"
func (s *Config) GetLogLevel() string {
	value := s.getString(LogLevel)
	if value != "Debug" && value != "Info" && value != "Warning" && value != "Error" && value != "Trace" {
		value = defaultLogLevel
	}

	return value
}

// GetLogAccess returns true if http requests should be logged to the terminal.
// HTTP requests are not logged to the log file. Defaults to true.
func (s *Config) GetLogAccess() bool {
	return s.getBoolDefault(LogAccess, defaultLogAccess)
}

// GetLogFileMaxSize returns the maximum size of the log file in megabytes for lumberjack to rotate
func (s *Config) GetLogFileMaxSize() int {
	value := s.getInt(LogFileMaxSize)
	if value < 0 {
		value = defaultLogFileMaxSize
	}

	return value
}

// Max allowed graphql upload size in megabytes
func (s *Config) GetMaxUploadSize() int64 {
	s.RLock()
	defer s.RUnlock()
	ret := int64(1024)

	v := s.forKey(MaxUploadSize)
	if v.Exists(MaxUploadSize) {
		ret = v.Int64(MaxUploadSize)
	}
	return ret << 20
}

// GetProxy returns the url of a http proxy to be used for all outgoing http calls.
func (s *Config) GetProxy() string {
	// Validate format
	reg := regexp.MustCompile(`^((?:socks5h?|https?):\/\/)(([\P{Cc}]+):([\P{Cc}]+)@)?(([a-zA-Z0-9][a-zA-Z0-9.-]*)(:[0-9]{1,5})?)`)
	proxy := s.getString(Proxy)
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
func (s *Config) GetNoProxy() string {
	// NoProxy does not require validation, it is validated by the native Go library sufficiently
	return s.getString(NoProxy)
}

// ActivatePublicAccessTripwire sets the security_tripwire_accessed_from_public_internet
// config field to the provided IP address to indicate that stash has been accessed
// from this public IP without authentication.
func (s *Config) ActivatePublicAccessTripwire(requestIP string) error {
	s.SetString(SecurityTripwireAccessedFromPublicInternet, requestIP)
	return s.Write()
}

func (s *Config) getPackageSources(key string) []*models.PackageSource {
	var sources []*models.PackageSource
	if err := s.unmarshalKey(key, &sources); err != nil {
		logger.Warnf("error in unmarshalkey: %v", err)
	}

	return sources
}

func (s *Config) GetPluginPackageSources() []*models.PackageSource {
	return s.getPackageSources(PluginPackageSources)
}

func (s *Config) GetScraperPackageSources() []*models.PackageSource {
	return s.getPackageSources(ScraperPackageSources)
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

func (s *Config) GetPluginPackagePathGetter() packagePathGetter {
	return packagePathGetter{
		getterFn: s.GetPluginPackageSources,
	}
}

func (s *Config) GetScraperPackagePathGetter() packagePathGetter {
	return packagePathGetter{
		getterFn: s.GetScraperPackageSources,
	}
}

func (s *Config) Validate() error {
	s.RLock()
	defer s.RUnlock()
	mandatoryPaths := []string{
		Database,
		Generated,
	}

	var missingFields []string

	for _, p := range mandatoryPaths {
		if !s.forKey(p).Exists(p) || s.forKey(p).String(p) == "" {
			missingFields = append(missingFields, p)
		}
	}

	if len(missingFields) > 0 {
		return MissingConfigError{
			missingFields: missingFields,
		}
	}

	if s.GetBlobsStorage() == BlobStorageTypeFilesystem && s.forKey(BlobsPath).String(BlobsPath) == "" {
		return MissingConfigError{
			missingFields: []string{BlobsPath},
		}
	}

	return nil
}

func (s *Config) setDefaultValues() {
	// read data before write lock scope
	defaultDatabaseFilePath := s.GetDefaultDatabaseFilePath()
	defaultScrapersPath := s.GetDefaultScrapersPath()
	defaultPluginsPath := s.GetDefaultPluginsPath()

	s.Lock()
	defer s.Unlock()

	// set the default host and port so that these are written to the config
	// file
	s.setDefault(Host, hostDefault)
	s.setDefault(Port, portDefault)

	s.setDefault(ParallelTasks, parallelTasksDefault)
	s.setDefault(SequentialScanning, SequentialScanningDefault)
	s.setDefault(PreviewSegmentDuration, previewSegmentDurationDefault)
	s.setDefault(PreviewSegments, previewSegmentsDefault)
	s.setDefault(PreviewExcludeStart, previewExcludeStartDefault)
	s.setDefault(PreviewExcludeEnd, previewExcludeEndDefault)
	s.setDefault(PreviewAudio, previewAudioDefault)
	s.setDefault(SoundOnPreview, false)

	s.setDefault(UseCustomSpriteInterval, UseCustomSpriteIntervalDefault)
	s.setDefault(SpriteInterval, SpriteIntervalDefault)
	s.setDefault(MinimumSprites, MinimumSpritesDefault)
	s.setDefault(MaximumSprites, MaximumSpritesDefault)
	s.setDefault(SpriteScreenshotSize, spriteScreenshotSizeDefault)

	s.setDefault(ThemeColor, DefaultThemeColor)

	s.setDefault(WriteImageThumbnails, writeImageThumbnailsDefault)
	s.setDefault(CreateImageClipsFromVideos, createImageClipsFromVideosDefault)

	s.setDefault(Database, defaultDatabaseFilePath)

	s.setDefault(dangerousAllowPublicWithoutAuth, dangerousAllowPublicWithoutAuthDefault)
	s.setDefault(SecurityTripwireAccessedFromPublicInternet, securityTripwireAccessedFromPublicInternetDefault)

	// Set generated to the metadata path for backwards compat
	s.setDefault(Generated, s.main.String(Metadata))

	s.setDefault(NoBrowser, NoBrowserDefault)
	s.setDefault(NotificationsEnabled, NotificationsEnabledDefault)
	s.setDefault(ShowOneTimeMovedNotification, ShowOneTimeMovedNotificationDefault)

	// Set default scrapers and plugins paths
	s.setDefault(ScrapersPath, defaultScrapersPath)
	s.setDefault(PluginsPath, defaultPluginsPath)

	// Set default gallery cover regex
	s.setDefault(GalleryCoverRegex, galleryCoverRegexDefault)

	// Set NoProxy default
	s.setDefault(NoProxy, noProxyDefault)

	// set default package sources
	s.setDefault(PluginPackageSources, []map[string]string{{
		"name":      sourceDefaultName,
		"url":       pluginPackageSourcesDefault,
		"localpath": sourceDefaultPath,
	}})
	s.setDefault(ScraperPackageSources, []map[string]string{{
		"name":      sourceDefaultName,
		"url":       scraperPackageSourcesDefault,
		"localpath": sourceDefaultPath,
	}})
}

// setExistingSystemDefaults sets config options that are new and unset in an existing install,
// but should have a separate default than for brand-new systems, to maintain behavior.
// The config file will not be written.
func (s *Config) setExistingSystemDefaults() {
	s.Lock()
	defer s.Unlock()
	if !s.isNewSystem {
		// Existing systems as of the introduction of auto-browser open should retain existing
		// behavior and not start the browser automatically.
		if !s.main.Exists(NoBrowser) {
			s.set(NoBrowser, true)
		}

		// Existing systems as of the introduction of the taskbar should inform users.
		if !s.main.Exists(ShowOneTimeMovedNotification) {
			s.set(ShowOneTimeMovedNotification, true)
		}
	}
}

// SetInitialConfig fills in missing required config fields. The config file will not be written.
func (s *Config) SetInitialConfig() error {
	// generate some api keys
	const apiKeyLength = 32

	if string(s.GetJWTSignKey()) == "" {
		signKey, err := hash.GenerateRandomKey(apiKeyLength)
		if err != nil {
			return fmt.Errorf("error generating JWTSignKey: %w", err)
		}
		s.SetString(JWTSignKey, signKey)
	}

	if string(s.GetSessionStoreKey()) == "" {
		sessionStoreKey, err := hash.GenerateRandomKey(apiKeyLength)
		if err != nil {
			return fmt.Errorf("error generating session store key: %w", err)
		}
		s.SetString(SessionStoreKey, sessionStoreKey)
	}

	s.setDefaultValues()

	return nil
}

func (s *Config) FinalizeSetup() {
	s.isNewSystem = false
	// i.configUpdates <- 0
}
