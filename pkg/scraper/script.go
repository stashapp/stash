package scraper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	stashJson "github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/python"
)

// inputs for scrapers

type fingerprintInput struct {
	Type        string `json:"type,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

type fileInput struct {
	ID      string             `json:"id"`
	ZipFile *fileInput         `json:"zip_file,omitempty"`
	ModTime stashJson.JSONTime `json:"mod_time"`

	Path string `json:"path,omitempty"`

	Fingerprints []fingerprintInput `json:"fingerprints,omitempty"`
	Size         int64              `json:"size,omitempty"`
}

type videoFileInput struct {
	fileInput
	Format     string  `json:"format,omitempty"`
	Width      int     `json:"width,omitempty"`
	Height     int     `json:"height,omitempty"`
	Duration   float64 `json:"duration,omitempty"`
	VideoCodec string  `json:"video_codec,omitempty"`
	AudioCodec string  `json:"audio_codec,omitempty"`
	FrameRate  float64 `json:"frame_rate,omitempty"`
	BitRate    int64   `json:"bitrate,omitempty"`

	Interactive      bool `json:"interactive,omitempty"`
	InteractiveSpeed *int `json:"interactive_speed,omitempty"`
}

// sceneInput is the input passed to the scraper for an existing scene
type sceneInput struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Code  string `json:"code,omitempty"`

	// deprecated - use urls instead
	URL  *string  `json:"url"`
	URLs []string `json:"urls"`

	// don't use omitempty for these to maintain backwards compatibility
	Date    *string `json:"date"`
	Details string  `json:"details"`

	Director string `json:"director,omitempty"`

	Files []videoFileInput `json:"files,omitempty"`
}

func fileInputFromFile(f models.BaseFile) fileInput {
	b := f.Base()
	var z *fileInput
	if b.ZipFile != nil {
		zz := fileInputFromFile(*b.ZipFile.Base())
		z = &zz
	}

	ret := fileInput{
		ID:      f.ID.String(),
		ZipFile: z,
		ModTime: stashJson.JSONTime{Time: f.ModTime},
		Path:    f.Path,
		Size:    f.Size,
	}

	for _, fp := range f.Fingerprints {
		ret.Fingerprints = append(ret.Fingerprints, fingerprintInput{
			Type:        fp.Type,
			Fingerprint: fp.Value(),
		})
	}

	return ret
}

func videoFileInputFromVideoFile(vf *models.VideoFile) videoFileInput {
	return videoFileInput{
		fileInput:        fileInputFromFile(*vf.Base()),
		Format:           vf.Format,
		Width:            vf.Width,
		Height:           vf.Height,
		Duration:         vf.Duration,
		VideoCodec:       vf.VideoCodec,
		AudioCodec:       vf.AudioCodec,
		FrameRate:        vf.FrameRate,
		BitRate:          vf.BitRate,
		Interactive:      vf.Interactive,
		InteractiveSpeed: vf.InteractiveSpeed,
	}
}

func sceneInputFromScene(scene *models.Scene) sceneInput {
	dateToStringPtr := func(s *models.Date) *string {
		if s != nil {
			v := s.String()
			return &v
		}

		return nil
	}

	// fallback to file basename if title is empty
	title := scene.GetTitle()

	var url *string
	urls := scene.URLs.List()
	if len(urls) > 0 {
		url = &urls[0]
	}

	ret := sceneInput{
		ID:      strconv.Itoa(scene.ID),
		Title:   title,
		Details: scene.Details,
		// include deprecated URL for now
		URL:      url,
		URLs:     urls,
		Date:     dateToStringPtr(scene.Date),
		Code:     scene.Code,
		Director: scene.Director,
	}

	for _, f := range scene.Files.List() {
		vf := videoFileInputFromVideoFile(f)
		ret.Files = append(ret.Files, vf)
	}

	return ret
}

type galleryInput struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Urls    []string `json:"urls"`
	Date    *string  `json:"date"`
	Details string   `json:"details"`

	Code         string `json:"code,omitempty"`
	Photographer string `json:"photographer,omitempty"`

	Files []fileInput `json:"files,omitempty"`

	// deprecated
	URL *string `json:"url"`
}

func galleryInputFromGallery(gallery *models.Gallery) galleryInput {
	dateToStringPtr := func(s *models.Date) *string {
		if s != nil {
			v := s.String()
			return &v
		}

		return nil
	}

	// fallback to file basename if title is empty
	title := gallery.GetTitle()

	var url *string
	urls := gallery.URLs.List()
	if len(urls) > 0 {
		url = &urls[0]
	}

	ret := galleryInput{
		ID:           strconv.Itoa(gallery.ID),
		Title:        title,
		Details:      gallery.Details,
		URL:          url,
		Urls:         urls,
		Date:         dateToStringPtr(gallery.Date),
		Code:         gallery.Code,
		Photographer: gallery.Photographer,
	}

	for _, f := range gallery.Files.List() {
		fi := fileInputFromFile(*f.Base())
		ret.Files = append(ret.Files, fi)
	}

	return ret
}

var ErrScraperScript = errors.New("scraper script error")

type scriptScraper struct {
	scraper      scraperTypeConfig
	config       config
	globalConfig GlobalConfig
}

func newScriptScraper(scraper scraperTypeConfig, config config, globalConfig GlobalConfig) *scriptScraper {
	return &scriptScraper{
		scraper:      scraper,
		config:       config,
		globalConfig: globalConfig,
	}
}

func (s *scriptScraper) runScraperScript(ctx context.Context, inString string, out interface{}) error {
	command := s.scraper.Script

	var cmd *exec.Cmd
	if python.IsPythonCommand(command[0]) {
		pythonPath := s.globalConfig.GetPythonPath()
		p, err := python.Resolve(pythonPath)

		if err != nil {
			logger.Warnf("%s", err)
		} else {
			cmd = p.Command(ctx, command[1:])
			envVariable, _ := filepath.Abs(filepath.Dir(filepath.Dir(s.config.path)))
			python.AppendPythonPath(cmd, envVariable)
		}
	}

	if cmd == nil {
		// if could not find python, just use the command args as-is
		cmd = stashExec.CommandContext(ctx, command[0], command[1:]...)
	}

	cmd.Dir = filepath.Dir(s.config.path)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()

		if n, err := io.WriteString(stdin, inString); err != nil {
			logger.Warnf("failure to write full input to script (wrote %v bytes out of %v): %v", n, len(inString), err)
		}
	}()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("Scraper stderr not available: " + err.Error())
	}

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Error("Scraper stdout not available: " + err.Error())
	}

	if err = cmd.Start(); err != nil {
		logger.Error("Error running scraper script: " + err.Error())
		return errors.New("error running scraper script")
	}

	go handleScraperStderr(s.config.Name, stderr)

	logger.Debugf("Scraper script <%s> started", strings.Join(cmd.Args, " "))

	// TODO - add a timeout here
	// Make a copy of stdout here. This allows us to decode it twice.
	var sb strings.Builder
	tr := io.TeeReader(stdout, &sb)

	// First, perform a decode where unknown fields are disallowed.
	d := json.NewDecoder(tr)
	d.DisallowUnknownFields()
	strictErr := d.Decode(out)

	if strictErr != nil {
		// The decode failed for some reason, use the built string
		// and allow unknown fields in the decode.
		s := sb.String()
		lenientErr := json.NewDecoder(strings.NewReader(s)).Decode(out)
		if lenientErr != nil {
			// The error is genuine, so return it
			logger.Errorf("could not unmarshal json from script output: %v", lenientErr)
			return fmt.Errorf("could not unmarshal json from script output: %w", lenientErr)
		}

		// Lenient decode succeeded, print a warning, but use the decode
		logger.Warnf("reading script result: %v", strictErr)
	}

	err = cmd.Wait()
	logger.Debugf("Scraper script finished")

	if err != nil {
		return fmt.Errorf("%w: %v", ErrScraperScript, err)
	}

	return nil
}

func (s *scriptScraper) scrapeByName(ctx context.Context, name string, ty ScrapeContentType) ([]ScrapedContent, error) {
	input := `{"name": "` + name + `"}`

	var ret []ScrapedContent
	var err error
	switch ty {
	case ScrapeContentTypePerformer:
		var performers []models.ScrapedPerformer
		err = s.runScraperScript(ctx, input, &performers)
		if err == nil {
			for _, p := range performers {
				v := p
				ret = append(ret, &v)
			}
		}
	case ScrapeContentTypeScene:
		var scenes []models.ScrapedScene
		err = s.runScraperScript(ctx, input, &scenes)
		if err == nil {
			for _, s := range scenes {
				v := s
				ret = append(ret, &v)
			}
		}
	default:
		return nil, ErrNotSupported
	}

	return ret, err
}

func (s *scriptScraper) scrapeByFragment(ctx context.Context, input Input) (ScrapedContent, error) {
	var inString []byte
	var err error
	var ty ScrapeContentType
	switch {
	case input.Performer != nil:
		inString, err = json.Marshal(*input.Performer)
		ty = ScrapeContentTypePerformer
	case input.Gallery != nil:
		inString, err = json.Marshal(*input.Gallery)
		ty = ScrapeContentTypeGallery
	case input.Scene != nil:
		inString, err = json.Marshal(*input.Scene)
		ty = ScrapeContentTypeScene
	}

	if err != nil {
		return nil, err
	}

	return s.scrape(ctx, string(inString), ty)
}

func (s *scriptScraper) scrapeByURL(ctx context.Context, url string, ty ScrapeContentType) (ScrapedContent, error) {
	return s.scrape(ctx, `{"url": "`+url+`"}`, ty)
}

func (s *scriptScraper) scrape(ctx context.Context, input string, ty ScrapeContentType) (ScrapedContent, error) {
	switch ty {
	case ScrapeContentTypePerformer:
		var performer *models.ScrapedPerformer
		err := s.runScraperScript(ctx, input, &performer)
		return performer, err
	case ScrapeContentTypeGallery:
		var gallery *models.ScrapedGallery
		err := s.runScraperScript(ctx, input, &gallery)
		return gallery, err
	case ScrapeContentTypeScene:
		var scene *models.ScrapedScene
		err := s.runScraperScript(ctx, input, &scene)
		return scene, err
	case ScrapeContentTypeMovie, ScrapeContentTypeGroup:
		var movie *models.ScrapedMovie
		err := s.runScraperScript(ctx, input, &movie)
		return movie, err
	case ScrapeContentTypeImage:
		var image *models.ScrapedImage
		err := s.runScraperScript(ctx, input, &image)
		return image, err
	}

	return nil, ErrNotSupported
}

func (s *scriptScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error) {
	inString, err := json.Marshal(sceneInputFromScene(scene))

	if err != nil {
		return nil, err
	}

	var ret *models.ScrapedScene

	err = s.runScraperScript(ctx, string(inString), &ret)

	return ret, err
}

func (s *scriptScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error) {
	inString, err := json.Marshal(galleryInputFromGallery(gallery))

	if err != nil {
		return nil, err
	}

	var ret *models.ScrapedGallery

	err = s.runScraperScript(ctx, string(inString), &ret)

	return ret, err
}

func (s *scriptScraper) scrapeImageByImage(ctx context.Context, image *models.Image) (*models.ScrapedImage, error) {
	inString, err := json.Marshal(imageToUpdateInput(image))

	if err != nil {
		return nil, err
	}

	var ret *models.ScrapedImage

	err = s.runScraperScript(ctx, string(inString), &ret)

	return ret, err
}

func handleScraperStderr(name string, scraperOutputReader io.ReadCloser) {
	const scraperPrefix = "[Scrape / %s] "

	lgr := logger.PluginLogger{
		Logger:          logger.Logger,
		Prefix:          fmt.Sprintf(scraperPrefix, name),
		DefaultLogLevel: &logger.ErrorLevel,
	}
	lgr.ReadLogMessages(scraperOutputReader)
}
