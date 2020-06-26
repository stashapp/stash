# Metadata Scraping

Stash supports scraping of performer and scene details.

Stash includes a freeones.xxx performer scraper built in.

# Adding custom scrapers

By default, Stash looks for scraper configurations in the `scrapers` sub-directory of the directory where the stash `config.yml` is read. This will either be the `$HOME/.stash` directory or the current working directory.

Custom scrapers are added by adding configuration yaml files (format: `scrapername.yml`) to the `scrapers` directory.

After scrapers are added, removed or edited while stash is running, they can be reloaded by clicking the `Scrape With...` button in New/Edit Performer or Scene page and clicking `Reload Scrapers`.

# Using custom scrapers

Scrapers support a number of different scraping types.

Performer details can be scraped from the new/edit Performer page in two different ways:

* click on the `Scrape With...` button and select the scraper to scrape with. You will be presented with a search dialog to search for the performer by their name
* enter the URL containing the Performer's details in the URL field. If the URL matches a pattern known to one of the scrapers, then a button will appear to scrape the details.

Scene details can be scraped using URL as above, or via the `Scrape With...` button, which scrapes using the current scene metadata.

# Community Scrapers
The stash community maintains a number of custom scraper configuration files that can be found [here](https://github.com/stashapp/CommunityScrapers).

# Scraper configuration file format

## Basic scraper configuration file structure

```
name: <site>
performerByName:
  <single scraper config>
performerByFragment:
  <single scraper config>
performerByURL:
  <multiple scraper URL configs>
sceneByFragment:
  <single scraper config>
sceneByURL:
  <multiple scraper URL configs>
<other configurations>
```

`name` is mandatory, all other top-level fields are optional. The inclusion of each top-level field determines what capabilities the scraper has.

A scraper configuration in any of the top-level fields must at least have an `action` field. The other fields are required based on the value of the `action` field.

The scraping types and their required fields are outlined in the following table:

| Behaviour | Required configuration |
|-----------|------------------------|
| Scraper in `Scrape...` dropdown button in Performer Edit page | Valid `performerByName` and `performerByFragment` configurations. |
| Scrape performer from URL | Valid `performerByURL` configuration with matching URL. |
| Scraper in `Scrape...` dropdown button in Scene Edit page | Valid `sceneByFragment` configuration. |
| Scrape scene from URL | Valid `sceneByURL` configuration with matching URL. |

URL-based scraping accepts multiple scrape configurations, and each configuration requires a `url` field. stash iterates through these configurations, attempting to match the entered URL against the `url` fields in the configuration. It executes the first scraping configuration where the entered URL contains the value of the `url` field. 

## Scraper Actions

### Script

Executes a script to perform the scrape. The `script` field is required for this action and accepts a list of string arguments. For example:

```
action: script
script:
  - python
  - iafdScrape.py
  - query
```

This configuration would execute `python iafdScrape.py query`.

Stash sends data to the script process's `stdin` stream and expects the output to be streamed to the `stdout` stream. Any errors and progress messages should be output to `stderr`.

The script is sent input and expects output based on the scraping type, as detailed in the following table:

| Scrape type | Input | Output |
|-------------|-------|--------|
| `performerByName` | `{"name": "<performer query string>"}` | Array of JSON-encoded performer fragments (including at least `name`) |
| `performerByFragment` | JSON-encoded performer fragment | JSON-encoded performer fragment |
| `performerByURL` | `{"url": "<url>"}` | JSON-encoded performer fragment |
| `sceneByFragment` | JSON-encoded scene fragment | JSON-encoded scene fragment |
| `sceneByURL` | `{"url": "<url>"}` | JSON-encoded scene fragment |

For `performerByName`, only `name` is required in the returned performer fragments. One entire object is sent back to `performerByFragment` to scrape a specific performer, so the other fields may be included to assist in scraping a performer. For example, the `url` field may be filled in for the specific performer page, then `performerByFragment` can extract by using its value.

As an example, the following python code snippet can be used to scrape a performer:

```
import json
import sys
import string

def readJSONInput():
	input = sys.stdin.read()
	return json.loads(input)

def searchPerformer(name):
    # perform scraping here - using name for the query

    # fill in the output
    ret = []
    
    # example shown for a single found performer 
    p = {}
    p['name'] = "some name"
    p['url'] = "performer url"
    ret.append(p)
    
    return ret

def scrapePerformer(input):
    ret = []
    # get the url from the input
    url = input['url']
    return scrapePerformerURL(url)

def debugPrint(t):
    sys.stderr.write(t + "\n")

def scrapePerformerURL(url):
    debugPrint("Reading url...")
    debugPrint("Parsing html...")
    
    # parse html

    # fill in performer details - single object
    ret = {}

    ret['name'] = "fred"
    ret['aliases'] = "freddy"
    ret['ethnicity'] = ""
    # and so on

    return ret

# read the input 
i = readJSONInput()

if sys.argv[1] == "query":
    ret = searchPerformer(i['name'])
    print(json.dumps(ret))
elif sys.argv[1] == "scrape":
    ret = scrapePerformer(i)
    print(json.dumps(ret))
elif sys.argv[1] == "scrapeURL":
    ret = scrapePerformerURL(i['url'])
    print(json.dumps(ret))
```

### scrapeXPath

This action scrapes a web page using an xpath configuration to parse. This action is valid for `performerByName`, `performerByURL` and `sceneByURL` only.

This action requires that the top-level `xPathScrapers` configuration is populated. The `scraper` field is required and must match the name of a scraper name configured in `xPathScrapers`. For example:

```
sceneByURL:
- action: scrapeXPath
  url: 
    - pornhub.com/view_video.php
  scraper: sceneScraper
```

The above configuration requires that `sceneScraper` exists in the `xPathScrapers` configuration.

#### Use with `performerByName`

For `performerByName`, the `queryURL` field must be present also. This field is used to perform a search query URL for performer names. The placeholder string sequence `{}` is replaced with the performer name search string. For the subsequent performer scrape to work, the `URL` field must be filled in with the URL of the performer page that matches a URL given in a `performerByURL` scraping configuration. For example:

```
name: Boobpedia
performerByName:
  action: scrapeXPath
  queryURL: http://www.boobpedia.com/wiki/index.php?title=Special%3ASearch&search={}&fulltext=Search
  scraper: performerSearch
performerByURL:
  - action: scrapeXPath
    url: 
      - boobpedia.com/boobs/
    scraper: performerScraper
xPathScrapers:
  performerSearch:
    performer:
      Name: # name element
      URL: # URL element that matches the boobpedia.com/boobs/ URL above
  performerScraper:
    # ... performer scraper details ...
```

#### XPath scrapers configuration

The top-level `xPathScrapers` field contains xpath scraping configurations, freely named. The scraping configuration may contain a `common` field, and must contain `performer` or `scene` depending on the scraping type it is configured for. 

Within the `performer`/`scene` field are key/value pairs corresponding to the golang fields (see below) on the performer/scene object. These fields are case-sensitive. 

The values of these may be either a simple xpath value, which tells the system where to get the value of the field from, or a more advanced configuration (see below). For example:

```
performer:
  Name: //h1[@itemprop="name"]
```

This will set the `Name` attribute of the returned performer to the text content of the element that matches `<h1 itemprop="name">...`.

The value may also be a sub-object, indicating that post-processing is required. If it is a sub-object, then the xpath must be set to the `selector` key of the sub-object. For example, using the same xpath as above:

```
performer:
  Name: 
    selector: //h1[@itemprop="name"]
    # post-processing config values
```

##### Common fragments

The `common` field is used to configure xpath fragments that can be referenced in the xpath strings. These are key-value pairs where the key is the string to reference the fragment, and the value is the string that the fragment will be replaced with. For example:

```
common:
  $infoPiece: //div[@class="infoPiece"]/span
performer:
  Measurements: $infoPiece[text() = 'Measurements:']/../span[@class="smallInfo"]  
```

The `Measurements` xpath string will replace `$infoPiece` with `//div[@class="infoPiece"]/span`, resulting in: `//div[@class="infoPiece"]/span[text() = 'Measurements:']/../span[@class="smallInfo"]`.

##### Post-processing options

The following post-processing keys are available:
* `concat`: if an xpath matches multiple elements, and `concat` is present, then all of the elements will be concatenated together
* `replace`: contains an array of sub-objects. Each sub-object must have a `regex` and `with` field. The `regex` field is the regex pattern to replace, and `with` is the string to replace it with. `$` is used to reference capture groups - `` is the first capture group, `` the second and so on. Replacements are performed in order of the array.
* `subScraper`: if present, the sub-scraper will be executed after all other post-processes are complete and before parseDate. It then takes the value and performs an http request, using the value as the URL. Within the `subScraper` config is a nested scraping configuration. This allows you to traverse to other webpages to get the attribute value you are after. For more info and examples have a look at [#370](https://github.com/stashapp/stash/pull/370), [#606](https://github.com/stashapp/stash/pull/606)
* `parseDate`: if present, the value is the date format using go's reference date (2006-01-02). For example, if an example date was `14-Mar-2003`, then the date format would be `02-Jan-2006`. See the [time.Parse documentation](https://golang.org/pkg/time/#Parse) for details. When present, the scraper will convert the input string into a date, then convert it to the string format used by stash (`YYYY-MM-DD`).
* `split`: Its the inverse of `concat`. Splits a string to more elements using the separator given. For more info and examples have a look at PR [#579](https://github.com/stashapp/stash/pull/579)

Post-processing is done in order of the fields above - `concat`, `regex`, `subscraper`, `parseDate` and then `split`.

##### Example

A performer and scene xpath scraper is shown as an example below:

```
name: Pornhub
performerByURL:
  - action: scrapeXPath
    url: 
      - pornhub.com
    scraper: performerScraper
sceneByURL:
  - action: scrapeXPath
    url: 
      - pornhub.com/view_video.php
    scraper: sceneScraper
xPathScrapers:
  performerScraper:
    common:
      $infoPiece: //div[@class="infoPiece"]/span
    performer:
      Name: //h1[@itemprop="name"]
      Birthdate: 
        selector: //span[@itemprop="birthDate"]
        parseDate: Jan 2, 2006
      Twitter: //span[text() = 'Twitter']/../@href
      Instagram: //span[text() = 'Instagram']/../@href
      Measurements: $infoPiece[text() = 'Measurements:']/../span[@class="smallInfo"]
      Height: 
        selector: $infoPiece[text() = 'Height:']/../span[@class="smallInfo"]
        replace: 
          - regex: .*\((\d+) cm\)
            with: $1
      Ethnicity: $infoPiece[text() = 'Ethnicity:']/../span[@class="smallInfo"]
      FakeTits: $infoPiece[text() = 'Fake Boobs:']/../span[@class="smallInfo"]
      Piercings: $infoPiece[text() = 'Piercings:']/../span[@class="smallInfo"]
      Tattoos: $infoPiece[text() = 'Tattoos:']/../span[@class="smallInfo"]
      CareerLength: 
        selector: $infoPiece[text() = 'Career Start and End:']/../span[@class="smallInfo"]
        replace:
          - regex: \s+to\s+
            with: "-"
  sceneScraper:
    common:
      $performer: //div[@class="pornstarsWrapper"]/a[@data-mxptype="Pornstar"]
      $studio: //div[@data-type="channel"]/a
    scene:
      Title: //div[@id="main-container"]/@data-video-title
      Tags: 
        Name: //div[@class="categoriesWrapper"]//a[not(@class="add-btn-small ")]
      Performers:
        Name: $performer/@data-mxptext
        URL: $performer/@href
      Studio:
        Name: $studio
        URL: $studio/@href    
```

See also [#333](https://github.com/stashapp/stash/pull/333) for more examples.

#### XPath resources:

- Test XPaths in Firefox: https://addons.mozilla.org/en-US/firefox/addon/try-xpath/
- XPath cheatsheet: https://devhints.io/xpath

#### Object fields
##### Performer

```
Name
Gender
URL
Twitter
Instagram
Birthdate
Ethnicity
Country
EyeColor
Height
Measurements
FakeTits
CareerLength
Tattoos
Piercings
Aliases
Image
```

*Note:*  - `Gender` must be one of `male`, `female`, `transgender_male`, `transgender_female` (case insensitive).

##### Scene
```
Title
Details
URL
Date
Image
Studio (see Studio Fields)
Movies (see Movie Fields)
Tags (see Tag fields)
Performers (list of Performer fields)
```
##### Studio
```
Name
URL
```

##### Tag
```
Name
```

##### Movie
```
Name
Aliases
Duration
Date
Rating
Director
Synopsis
URL
```

### Stash

A different stash server can be configured as a scraping source. This action applies only to `performerByName`, `performerByFragment`, and `sceneByFragment` types. This action requires that the top-level `stashServer` field is configured.

`stashServer` contains a single `url` field for the remote stash server. The username and password can be embedded in this string using `username:password@host`.

An example stash scrape configuration is below:

```
name: stash
performerByName:
  action: stash
performerByFragment:
  action: stash
sceneByFragment:
  - action: stash
stashServer:
  url: http://stashserver.com:9999
```

### Debugging support
To print the received html from a scraper request to the log file, add the following to your scraper yml file:
```
debug:
  printHTML: true
```
