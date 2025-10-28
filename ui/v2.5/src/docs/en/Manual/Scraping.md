# Metadata Scraping

Stash supports scraping of metadata from various external sources.

## Scraper Types

| Type | Description |
|---|:---|
| Fragment | Uses existing metadata for an Item and match it to a result from a metadata source. |
| Search/By Name | Uses a provided query string to search a metadata source for a list of matches for the user to pick from. |
| URL | Extracts metadata from a given URL. |

## Supported Scrapers

|   | Fragment | Search | URL |
|---|:---:|:---:|:---:|
| gallery | ✔️ | | ✔️ |
| image | ✔️ | | ✔️ |
| group | | | ✔️ |
| performer | | ✔️ | ✔️ |
| scene | ✔️  | ✔️ | ✔️ |

## Included Scrapers

Stash provides the following built-in scrapers:

| Scraper | Description |
|---|--|
| Freeones | `search` Performer scraper for freeones.xxx. |
| Auto Tag | Scene `fragment` scraper that matches existing performers, studio and tags using the filename. |

## Managing Scrapers

Scrapers can be installed and managed from the `Settings > Metadata Providers` page. 

Scrapers are installed using the `Available Scrapers` section. This section allows configuring sources from which to install scrapers. The `Community (stable)` source is configured by default. This source contains scrapers for the current _stable_ version of stash.

These are the scraper sources maintained by the stashapp organisation:

| Name | Source URL | Recommended Local Path | Notes |
|------|-----------|------------------------|-------|
| Community (stable) | `https://stashapp.github.io/CommunityScrapers/stable/index.yml` | `stable` | For the current stable version of stash. |
| Community (develop) | `https://stashapp.github.io/CommunityScrapers/develop/index.yml` | `develop` | For the develop version of stash. |

Installed scrapers can be updated or uninstalled from the `Installed Scrapers` section.

### Source URLs

The source URL must return a yaml file containing all the available packages for the source. An example source yaml file looks like the following:

```
- id: <package id>
  name: <package name>
  version: <version>
  date: <date>
  requires:
  - <ids of packages required by this package (optional)>
  - ...
  path: <path to package zip file>
  sha256: <sha256 of zip>
  metadata:
    <optional key/value pairs for extra information>
- ...
```

Path can be a relative path to the zip file or an external URL.

## Adding Scrapers manually

By default, Stash looks for scraper configurations in the `scrapers` sub-directory of the directory where the stash `config.yml` is read. This will either be the `$HOME/.stash` directory or the current working directory.

Scrapers are added manually by placing yaml configuration files (format: `scrapername.yml`) in the `scrapers` directory.

> **⚠️ Note:** Some scrapers may require more than just the yaml file, consult the individual scraper documentation

After the yaml files are added, removed or edited while stash is running, they can be reloaded going to `Settings > Metadata Providers > Scrapers` and clicking `Reload Scrapers`.
  
## Using Scrapers

#### Fragment Scraper
Click on the `Scrape With...` button in the `edit` tab of an item, then select the scraper you wish to use.

#### Search Scraper
Click on the 🔍 button in the `edit` tab of an item. You will be presented with a search dialog with a pre-populated query to search for, after searching you will be presented with a list of results to pick from

#### URL Scraper
Enter the URL in the `edit` tab of an Item. If a scraper is installed that supports that url, then a button will appear to scrape the metadata.

## Tagger View

The Tagger view is accessed from the scenes page. It allows the user to run scrapers on all items on the current page. The Tagger presents the user with potential matches for an item from a selected stash-box instance or metadata source if supported. The user needs to select the correct metadata information to save. 

When used in combination with stash-box, the user can optionally submit scene fingerprints to contribute to a stash-box instance. A scene fingerprint consists of any generated hashes (`phash`, `oshash`, `md5`) and the scene duration. Fingerprint submissions are associated with your stash-box account. Submitting fingerprints assists others in matching their files, because stash-box returns a count of matching user submitted fingerprints with every potential match.

| | Has Tagger | Source Selection |
|---|:---:|:---:|
| gallery | | |
| group | | |
| performer | ✔️ | |
| scene | ✔️ | ✔️ |


## Identify Task

This task iterates through your Scenes and attempts to identify the scene using a selection of scraping sources. This task can be found under `Settings -> Tasks -> "Identify..." (Button)`. For more information see the [Tasks > Identify](/help/Identify.md) page.
