package scraper

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
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

	// TODO - add a timeout here
	decodeErr := json.NewDecoder(stdout).Decode(out)

	stderrData, _ := ioutil.ReadAll(stderr)
	stderrString := string(stderrData)

	err = cmd.Wait()

	if err != nil {
		// error message should be in the stderr stream
		logger.Errorf("scraper error when running command <%s>: %s", strings.Join(cmd.Args, " "), stderrString)
		return errors.New("Error running scraper script")
	}

	if decodeErr != nil {
		logger.Errorf("error decoding performer from scraper data: %s", err.Error())
		return errors.New("Error decoding performer from scraper script")
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

func (s *scriptScraper) scrapeSceneByURL(url string) (*models.ScrapedScene, error) {
	inString := `{"url": "` + url + `"}`

	var ret models.ScrapedScene

	err := s.runScraperScript(string(inString), &ret)

	return &ret, err
}
