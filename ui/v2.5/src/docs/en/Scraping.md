# Metadata Scraping

Stash supports scraping of metadata from various external sources

## Terminology

| Term | Description |
|---|:---|
| Fragment Scraper | This scraper will attempt to use all existing metadata for an Item and match it to a result from a metadata source automatically |
| Search/By Name Scraper | This type of scraper will use the current name of the Item to search a given metadata source for a list of matches for the user to pick from|
| URL Scraper | This is a scraper that attempts to extract metadata from a given URL |

## Supported Scrapers

|   | Fragment | Search | URL |
|---|:---:|:---:|:---:|
| gallery | :heavy_check_mark: | | :heavy_check_mark: |
| movie | | | :heavy_check_mark: |
| performer | | :heavy_check_mark: | :heavy_check_mark: |
| scene | :heavy_check_mark:  | :heavy_check_mark: | :heavy_check_mark: |

# Scraper Operation

## Included Scrapers

Stash has a built-in performer `search` scraper for freeones.xxx.

## Adding Scrapers


By default, Stash looks for scraper configurations in the `scrapers` sub-directory of the directory where the stash `config.yml` is read. This will either be the `$HOME/.stash` directory or the current working directory.

Scrapers are added by placing yaml configuration files (format: `scrapername.yml`) in the `scrapers` directory.

> **⚠️ Note:** Some scrapers may require more than just the yaml file, consult the individual scraper documentation

After the yaml files are added, removed or edited while stash is running, they can be reloaded going to `Settings > Scrapers` and clicking `Reload Scrapers`.

The stash community maintains a number of custom scraper configuration files that can be found [here](https://github.com/stashapp/CommunityScrapers)
  
## Using Scrapers

#### Fragment Scraper
Click on the `Scrape With...` button in the `edit` tab of an item, then select the scraper you wish to use.

#### Search Scraper
Click on the 🔍 button in the `edit` tab of an item, You will be presented with a search dialog with a pre-populated query to search for, after searching you will be presented with a list of results to pick from

#### URL Scraper
Enter the URL in the `edit` tab of an Item, If a scraper is installed that supports that url then a button will appear to scrape the metadata.

## Tagger View

The Tagger refers to a specific view for items in stash stash that allows the user to run scrapers on those items, a page at a time, the tagger will present the user with potential matches for an item from a stash-box instance or from a selected metadata source if supported. The user is needed to select and save the correct metadata information to stash. 

When used in combination with stash-box the user can optionally submit fingerprints for scenes to contribute to a stash-box instance. Doing so will submit generated hashes (`phash`, `oshash`, `md5`) and the duration of the scene to assist others in matching their files based off these fingerprints. These are the only values stash submits to a stash-box instance, fingerprint submissions are anonymous.

| | Has Tagger | Source Selection |
|---|:---:|:---:|
| gallery | | |
| movie | | |
| performer | :heavy_check_mark: | |
| scene | :heavy_check_mark: | :heavy_check_mark: |


## Identify Task
The Identify task will automatically run a scraper against a group of files/folders and set the corresponding item metadata without user input for each item. This task can be found under `Settings -> Tasks -> "Identify..." (Button)` for more information see `Tasks > Identify` in this manual