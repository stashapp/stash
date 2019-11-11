package scraper

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
)

type ScraperMethod string

const (
	ScraperMethodScript ScraperMethod = "SCRIPT"
)

var AllScraperMethod = []ScraperMethod{
	ScraperMethodScript,
}

func (e ScraperMethod) IsValid() bool {
	switch e {
	case ScraperMethodScript:
		return true
	}
	return false
}

type scraperConfig struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Type              models.ScraperType `json:"type"`
	Method            ScraperMethod      `json:"method"`
	GetPerformerNames []string           `json:"get_performer_names"`
	GetPerformer      []string           `json:"get_performer"`
}

func (c scraperConfig) toScraper() *models.Scraper {
	ret := models.Scraper{
		ID:   c.ID,
		Name: c.Name,
		Type: c.Type,
	}

	return &ret
}

func (c scraperConfig) scrapePerformerNames(name string) ([]string, error) {
	if c.Method == ScraperMethodScript {
		return c.scrapePerformerNamesScript(name)
	}

	return nil, nil
}

func (c scraperConfig) scrapePerformer(name string) (*models.ScrapedPerformer, error) {
	if c.Method == ScraperMethodScript {
		return c.scrapePerformerScript(name)
	}

	return nil, nil
}

func runScraperScript(command []string, inString string, out interface{}) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = config.GetScrapersPath()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()

		// TODO - encode search query in proper json
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
		return err
	}

	// TODO - add a timeout here
	if err := json.NewDecoder(stdout).Decode(out); err != nil {
		return err
	}

	stderrData, _ := ioutil.ReadAll(stderr)
	stderrString := string(stderrData)

	err = cmd.Wait()

	if err != nil {
		// error message should be in the stderr stream
		logger.Errorf("scraper error when running command <%s>: %s", strings.Join(cmd.Args, " "), stderrString)
		return err
	}

	return nil
}

func (c scraperConfig) scrapePerformerNamesScript(name string) ([]string, error) {
	// TODO - encode search query in proper json
	inString := `{"name": "` + name + `"}`

	var ret []string

	err := runScraperScript(c.GetPerformerNames, inString, &ret)

	return ret, err
}

func (c scraperConfig) scrapePerformerScript(name string) (*models.ScrapedPerformer, error) {
	// TODO - encode search query in proper json
	inString := `{"name": "` + name + `"}`

	var ret models.ScrapedPerformer

	err := runScraperScript(c.GetPerformer, inString, &ret)

	return &ret, err
}

var scrapers []scraperConfig

func loadScraper(path string) (*scraperConfig, error) {
	var scraper scraperConfig
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&scraper)
	if err != nil {
		return nil, err
	}

	// set id to the filename
	id := filepath.Base(path)
	id = id[:strings.LastIndex(id, ".")]
	scraper.ID = id

	return &scraper, nil
}

func loadScrapers() ([]scraperConfig, error) {
	if scrapers != nil {
		return scrapers, nil
	}

	path := config.GetScrapersPath()
	scrapers = make([]scraperConfig, 0)

	logger.Infof("Reading scraper configs from %s", path)
	scraperFiles, err := filepath.Glob(filepath.Join(path, "*.json"))

	if err != nil {
		logger.Errorf("Error reading scraper configs: %s", err.Error())
		return nil, err
	}

	for _, file := range scraperFiles {
		scraper, err := loadScraper(file)
		if err != nil {
			logger.Errorf("Error loading scraper %s: %s", file, err.Error())
		} else {
			scrapers = append(scrapers, *scraper)
		}
	}

	return scrapers, nil
}

func ListScrapers(scraperType models.ScraperType) ([]*models.Scraper, error) {
	// read scraper config files from the directory and cache
	scrapers, err := loadScrapers()

	if err != nil {
		return nil, err
	}

	var ret []*models.Scraper
	for _, s := range scrapers {
		// filter on type
		if s.Type == scraperType {
			ret = append(ret, s.toScraper())
		}
	}

	return ret, nil
}

func findPerformerScraper(scraperID string) *scraperConfig {
	// read scraper config files from the directory and cache
	loadScrapers()

	for _, s := range scrapers {
		if s.ID == scraperID {
			return &s
		}
	}

	return nil
}

func ScrapePerformerList(scraperID string, query string) ([]string, error) {
	// find scraper with the provided id
	s := findPerformerScraper(scraperID)
	if s != nil {
		return s.scrapePerformerNames(query)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

func ScrapePerformer(scraperID string, performerName string) (*models.ScrapedPerformer, error) {
	// find scraper with the provided id
	s := findPerformerScraper(scraperID)
	if s != nil {
		return s.scrapePerformer(performerName)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}
