# Metadata Scraping

Stash supports scraping of metadata from various external sources

# Terminology
  ### URL Scraper
  * A URL scraper is a scraper that attempts to extract metadata from a given URL
  ### Fragment Scraper 
  * This type of scraper will attempt to use all current metadata for an Item and match it to a result from a metadata source automatically
  ### Search/By Name Search
  * This type of scraper will use the current name of the Item to search a given metadata source for a list of matches
  * This scraper relies on the user picking the correct match from a list of results 

# Supported Scrapers:

|   | URL | Search | Fragment |
|---|:---:|:---:|:---:|
| gallery | :heavy_check_mark: | | :heavy_check_mark: |
| movie | :heavy_check_mark: | | |
| performer | :heavy_check_mark: | :heavy_check_mark: |   |
| scene | :heavy_check_mark: | :heavy_check_mark:| :heavy_check_mark: |

# Scraper Operation

## Included Scrapers

Stash has a built in performer `search` scraper for freeones.xxx.

## Adding Scrapers


By default, Stash looks for scraper configurations in the `scrapers` sub-directory of the directory where the stash `config.yml` is read. This will either be the `$HOME/.stash` directory or the current working directory.

Scrapers are added by placing yaml configuration files (format: `scrapername.yml`) in the `scrapers` directory.

> **⚠️ Note:** Some scrapers may require more than just the yaml file, consult the individual scraper documentation

After the yaml files are added, removed or edited while stash is running, they can be reloaded going to `Settings > Scrapers` and clicking `Reload Scrapers`.

The stash community maintains a number of custom scraper configuration files that can be found [here](https://github.com/stashapp/CommunityScrapers)
  
## Using Scrapers

URL Scraper
* Enter the URL into the `edit` tab of an Item, If a scraper is installed that supports that url then a button will appear to scrape the metadata.

Fragment Scraper
* click on the `Scrape With...` button and select the scraper you wish to use.

Search Scraper
* Click on the :mag: button, You will be presented with a search dialog with a pre-populated query to search for, after searching you will be presented with a list of results to pick from
