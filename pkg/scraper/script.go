package scraper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/python"
)

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
		var p *python.Python
		if pythonPath != "" {
			p = python.New(pythonPath)
		} else {
			p, _ = python.Resolve()
		}

		if p != nil {
			cmd = p.Command(ctx, command[1:])
		}

		// if could not find python, just use the command args as-is
	}

	if cmd == nil {
		cmd = stashExec.Command(command[0], command[1:]...)
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

func (s *scriptScraper) scrapeByName(ctx context.Context, name string, ty models.ScrapeContentType) ([]models.ScrapedContent, error) {
	input := `{"name": "` + name + `"}`

	var ret []models.ScrapedContent
	var err error
	switch ty {
	case models.ScrapeContentTypePerformer:
		var performers []models.ScrapedPerformer
		err = s.runScraperScript(ctx, input, &performers)
		if err == nil {
			for _, p := range performers {
				v := p
				ret = append(ret, &v)
			}
		}
	case models.ScrapeContentTypeScene:
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

func (s *scriptScraper) scrapeByFragment(ctx context.Context, input Input) (models.ScrapedContent, error) {
	var inString []byte
	var err error
	var ty models.ScrapeContentType
	switch {
	case input.Performer != nil:
		inString, err = json.Marshal(*input.Performer)
		ty = models.ScrapeContentTypePerformer
	case input.Gallery != nil:
		inString, err = json.Marshal(*input.Gallery)
		ty = models.ScrapeContentTypeGallery
	case input.Scene != nil:
		inString, err = json.Marshal(*input.Scene)
		ty = models.ScrapeContentTypeScene
	}

	if err != nil {
		return nil, err
	}

	return s.scrape(ctx, string(inString), ty)
}

func (s *scriptScraper) scrapeByURL(ctx context.Context, url string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	return s.scrape(ctx, `{"url": "`+url+`"}`, ty)
}

func (s *scriptScraper) scrape(ctx context.Context, input string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	switch ty {
	case models.ScrapeContentTypePerformer:
		var performer *models.ScrapedPerformer
		err := s.runScraperScript(ctx, input, &performer)
		return performer, err
	case models.ScrapeContentTypeGallery:
		var gallery *models.ScrapedGallery
		err := s.runScraperScript(ctx, input, &gallery)
		return gallery, err
	case models.ScrapeContentTypeScene:
		var scene *models.ScrapedScene
		err := s.runScraperScript(ctx, input, &scene)
		return scene, err
	case models.ScrapeContentTypeMovie:
		var movie *models.ScrapedMovie
		err := s.runScraperScript(ctx, input, &movie)
		return movie, err
	}

	return nil, ErrNotSupported
}

func (s *scriptScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error) {
	inString, err := json.Marshal(sceneToUpdateInput(scene))

	if err != nil {
		return nil, err
	}

	var ret *models.ScrapedScene

	err = s.runScraperScript(ctx, string(inString), &ret)

	return ret, err
}

func (s *scriptScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error) {
	inString, err := json.Marshal(galleryToUpdateInput(gallery))

	if err != nil {
		return nil, err
	}

	var ret *models.ScrapedGallery

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
