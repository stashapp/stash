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
	ScraperMethodScript  ScraperMethod = "SCRIPT"
	ScraperMethodBuiltin ScraperMethod = "BUILTIN"
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
	URLs              []string           `json:"urls"`
	GetPerformerNames []string           `json:"get_performer_names"`
	GetPerformer      []string           `json:"get_performer"`
	GetPerformerURL   []string           `json:"get_performer_url"`

	scrapePerformerNamesFunc func(c scraperConfig, name string) ([]*models.ScrapedPerformer, error)
	scrapePerformerFunc      func(c scraperConfig, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error)
	scrapePerformerURLFunc   func(c scraperConfig, url string) (*models.ScrapedPerformer, error)
}

func (c scraperConfig) toScraper() *models.Scraper {
	ret := models.Scraper{
		ID:   c.ID,
		Name: c.Name,
		Type: c.Type,
		Urls: c.URLs,
	}

	return &ret
}

func (c *scraperConfig) postDecode() {
	if c.Method == ScraperMethodScript {
		c.scrapePerformerNamesFunc = scrapePerformerNamesScript
		c.scrapePerformerFunc = scrapePerformerScript
	}
}

func (c scraperConfig) ScrapePerformerNames(name string) ([]*models.ScrapedPerformer, error) {
	return c.scrapePerformerNamesFunc(c, name)
}

func (c scraperConfig) ScrapePerformer(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	return c.scrapePerformerFunc(c, scrapedPerformer)
}

func (c scraperConfig) ScrapePerformerURL(url string) (*models.ScrapedPerformer, error) {
	return c.scrapePerformerURLFunc(c, url)
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

func scrapePerformerNamesScript(c scraperConfig, name string) ([]*models.ScrapedPerformer, error) {
	inString := `{"name": "` + name + `"}`

	var performers []models.ScrapedPerformer

	err := runScraperScript(c.GetPerformerNames, inString, &performers)

	// convert to pointers
	var ret []*models.ScrapedPerformer
	if err == nil {
		for i := 0; i < len(performers); i++ {
			ret = append(ret, &performers[i])
		}
	}

	return ret, err
}

func scrapePerformerScript(c scraperConfig, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	inString, err := json.Marshal(scrapedPerformer)

	if err != nil {
		return nil, err
	}

	var ret models.ScrapedPerformer

	err = runScraperScript(c.GetPerformer, string(inString), &ret)

	return &ret, err
}

func scrapePerformerURLScript(c scraperConfig, url string) (*models.ScrapedPerformer, error) {
	inString := `{"url": "` + url + `"}`

	var ret models.ScrapedPerformer

	err := runScraperScript(c.GetPerformerURL, string(inString), &ret)

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
	scraper.postDecode()

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

	// add built-in freeones scraper
	scrapers = append(scrapers, GetFreeonesScraper())

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

func findPerformerScraperURL(url string) *scraperConfig {
	// read scraper config files from the directory and cache
	loadScrapers()

	for _, s := range scrapers {
		for _, url := range s.URLs {
			if strings.Contains(url, url) {
				return &s
			}
		}
	}

	return nil
}

func ScrapePerformerList(scraperID string, query string) ([]*models.ScrapedPerformer, error) {
	// find scraper with the provided id
	s := findPerformerScraper(scraperID)
	if s != nil {
		return s.ScrapePerformerNames(query)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

func ScrapePerformer(scraperID string, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	// find scraper with the provided id
	s := findPerformerScraper(scraperID)
	if s != nil {
		return s.ScrapePerformer(scrapedPerformer)
	}

	return nil, errors.New("Scraper with ID " + scraperID + " not found")
}

func ScrapePerformerURL(url string) (*models.ScrapedPerformer, error) {
	// find scraper that matches the url given
	s := findPerformerScraperURL(url)
	if s != nil {
		return s.ScrapePerformerURL(url)
	}

	return nil, nil
}
