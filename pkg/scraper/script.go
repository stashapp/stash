package scraper

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common/log"
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

		io.WriteString(stdin, inString)
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
		return errors.New("Error running scraper script")
	}

	go handleScraperStderr(stderr)

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
		return errors.New("Error running scraper script")
	}

	return nil
}

func (s *scriptScraper) scrapePerformersByName(name string) ([]*models.ScrapedPerformer, error) {
	inString := `{"name": "` + name + `"}`

	var performers []models.ScrapedPerformer

	err := s.runScraperScript(inString, &performers)

	// convert to pointers
	var ret []*models.ScrapedPerformer
	if err == nil {
		for i := 0; i < len(performers); i++ {
			ret = append(ret, &performers[i])
		}
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

func (s *scriptScraper) scrapePerformerByURL(url string) (*models.ScrapedPerformer, error) {
	inString := `{"url": "` + url + `"}`

	var ret models.ScrapedPerformer

	err := s.runScraperScript(string(inString), &ret)

	return &ret, err
}

func (s *scriptScraper) scrapeSceneByFragment(scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
	inString, err := json.Marshal(scene)

	if err != nil {
		return nil, err
	}

	var ret models.ScrapedScene

	err = s.runScraperScript(string(inString), &ret)

	return &ret, err
}

func (s *scriptScraper) scrapeGalleryByFragment(gallery models.GalleryUpdateInput) (*models.ScrapedGallery, error) {
	inString, err := json.Marshal(gallery)

	if err != nil {
		return nil, err
	}

	var ret models.ScrapedGallery

	err = s.runScraperScript(string(inString), &ret)

	return &ret, err
}

func (s *scriptScraper) scrapeSceneByURL(url string) (*models.ScrapedScene, error) {
	inString := `{"url": "` + url + `"}`

	var ret models.ScrapedScene

	err := s.runScraperScript(string(inString), &ret)

	return &ret, err
}

func (s *scriptScraper) scrapeGalleryByURL(url string) (*models.ScrapedGallery, error) {
	inString := `{"url": "` + url + `"}`

	var ret models.ScrapedGallery

	err := s.runScraperScript(string(inString), &ret)

	return &ret, err
}

func (s *scriptScraper) scrapeMovieByURL(url string) (*models.ScrapedMovie, error) {
	inString := `{"url": "` + url + `"}`

	var ret models.ScrapedMovie

	err := s.runScraperScript(string(inString), &ret)

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

func handleStderrLine(line string, defaultLogLevel *log.Level) {
	level, l := log.DetectLogLevel(line)

	const scraperPrefix = "[Scrape] "
	// if no log level, just output to info
	if level == nil {
		if defaultLogLevel != nil {
			level = defaultLogLevel
		} else {
			level = &log.InfoLevel
		}
	}

	switch *level {
	case log.TraceLevel:
		logger.Trace(scraperPrefix, l)
	case log.DebugLevel:
		logger.Debug(scraperPrefix, l)
	case log.InfoLevel:
		logger.Info(scraperPrefix, l)
	case log.WarningLevel:
		logger.Warn(scraperPrefix, l)
	case log.ErrorLevel:
		logger.Error(scraperPrefix, l)
	}
}

func handleScraperOutput(scraperOutputReader io.ReadCloser, defaultLogLevel *log.Level) {
	// pipe plugin stderr to our logging
	scanner := bufio.NewScanner(scraperOutputReader)
	for scanner.Scan() {
		str := scanner.Text()
		if str != "" {
			handleStderrLine(str, defaultLogLevel)
		}
	}

	str := scanner.Text()
	if str != "" {
		handleStderrLine(str, defaultLogLevel)
	}

	scraperOutputReader.Close()
}

func handleScraperStderr(scraperOutputReader io.ReadCloser) {
	handleScraperOutput(scraperOutputReader, &log.ErrorLevel)
}
