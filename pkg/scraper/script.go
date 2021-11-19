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

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

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

func (s *scriptScraper) runScraperScript(inString string, out interface{}) error {
	command := s.scraper.Script

	if command[0] == "python" || command[0] == "python3" {
		executable, err := findPythonExecutable()
		if err == nil {
			command[0] = executable
		}
	}

	cmd := exec.Command(command[0], command[1:]...)
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
	decodeErr := json.NewDecoder(stdout).Decode(out)
	if decodeErr != nil {
		logger.Error("could not unmarshal json: " + decodeErr.Error())
		return errors.New("could not unmarshal json: " + decodeErr.Error())
	}

	err = cmd.Wait()
	logger.Debugf("Scraper script finished")

	if err != nil {
		return errors.New("error running scraper script")
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
		err = s.runScraperScript(input, &performers)
		if err == nil {
			for _, p := range performers {
				v := p
				ret = append(ret, &v)
			}
		}
	case models.ScrapeContentTypeScene:
		var scenes []models.ScrapedScene
		err = s.runScraperScript(input, &scenes)
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

func (s *scriptScraper) scrapePerformerByFragment(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	inString, err := json.Marshal(scrapedPerformer)

	if err != nil {
		return nil, err
	}

	var ret models.ScrapedPerformer

	err = s.runScraperScript(string(inString), &ret)

	return &ret, err
}

func (s *scriptScraper) scrapeByURL(ctx context.Context, url string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	input := `{"url": "` + url + `"}`

	switch ty {
	case models.ScrapeContentTypePerformer:
		var performer models.ScrapedPerformer
		err := s.runScraperScript(input, &performer)
		return &performer, err
	case models.ScrapeContentTypeGallery:
		var gallery models.ScrapedGallery
		err := s.runScraperScript(input, &gallery)
		return &gallery, err
	case models.ScrapeContentTypeScene:
		var scene models.ScrapedScene
		err := s.runScraperScript(input, &scene)
		return &scene, err
	case models.ScrapeContentTypeMovie:
		var movie models.ScrapedMovie
		err := s.runScraperScript(input, &movie)
		return &movie, err
	}

	return nil, ErrNotSupported
}

func (s *scriptScraper) scrapeSceneByScene(ctx context.Context, scene *models.Scene) (*models.ScrapedScene, error) {
	inString, err := json.Marshal(sceneToUpdateInput(scene))

	if err != nil {
		return nil, err
	}

	var ret models.ScrapedScene

	err = s.runScraperScript(string(inString), &ret)

	return &ret, err
}

func (s *scriptScraper) scrapeSceneByFragment(ctx context.Context, scene models.ScrapedSceneInput) (*models.ScrapedScene, error) {
	inString, err := json.Marshal(scene)

	if err != nil {
		return nil, err
	}

	var ret models.ScrapedScene

	err = s.runScraperScript(string(inString), &ret)

	return &ret, err
}

func (s *scriptScraper) scrapeGalleryByGallery(ctx context.Context, gallery *models.Gallery) (*models.ScrapedGallery, error) {
	inString, err := json.Marshal(galleryToUpdateInput(gallery))

	if err != nil {
		return nil, err
	}

	var ret models.ScrapedGallery

	err = s.runScraperScript(string(inString), &ret)

	return &ret, err
}

func (s *scriptScraper) scrapeGalleryByFragment(gallery models.ScrapedGalleryInput) (*models.ScrapedGallery, error) {
	inString, err := json.Marshal(gallery)

	if err != nil {
		return nil, err
	}

	var ret models.ScrapedGallery

	err = s.runScraperScript(string(inString), &ret)

	return &ret, err
}

func findPythonExecutable() (string, error) {
	_, err := exec.LookPath("python3")

	if err != nil {
		_, err = exec.LookPath("python")

		if err != nil {
			return "", err
		}

		return "python", nil
	}

	return "python3", nil
}

func handleScraperStderr(name string, scraperOutputReader io.ReadCloser) {
	const scraperPrefix = "[Scrape / %s] "

	lgr := logger.PluginLogger{
		Prefix:          fmt.Sprintf(scraperPrefix, name),
		DefaultLogLevel: &logger.ErrorLevel,
	}
	lgr.HandlePluginStdErr(scraperOutputReader)
}
