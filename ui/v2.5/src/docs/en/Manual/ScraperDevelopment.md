# Contributing Scrapers 

Scrapers can be contributed to the community by creating a PR in [this repository](https://github.com/stashapp/CommunityScrapers/pulls).

## Scraper configuration file format

```yaml
name: <site>
performerByName:
  <single scraper config>
performerByFragment:
  <single scraper config>
performerByURL:
  <multiple scraper URL configs>
sceneByName:
  <single scraper config>
sceneByQueryFragment:
  <single scraper config>
sceneByFragment:
  <single scraper config>
sceneByURL:
  <multiple scraper URL configs>
movieByURL:
  <multiple scraper URL configs>
galleryByFragment:
  <single scraper config>
galleryByURL:
  <multiple scraper URL configs>
<other configurations>
```

`name` is mandatory, all other top-level fields are optional. The inclusion of each top-level field determines what capabilities the scraper has.

A scraper configuration in any of the top-level fields must at least have an `action` field. The other fields are required based on the value of the `action` field.

The scraping types and their required fields are outlined in the following table:

| Behavior | Required configuration |
|-----------|------------------------|
| Scraper in `Scrape...` dropdown button in Performer Edit page | Valid `performerByName` and `performerByFragment` configurations. |
| Scrape performer from URL | Valid `performerByURL` configuration with matching URL. |
| Scraper in query dropdown button in Scene Edit page | Valid `sceneByName` and `sceneByQueryFragment` configurations. |
| Scraper in `Scrape...` dropdown button in Scene Edit page | Valid `sceneByFragment` configuration. |
| Scrape scene from URL | Valid `sceneByURL` configuration with matching URL. |
| Scrape movie from URL | Valid `movieByURL` configuration with matching URL. |
| Scraper in `Scrape...` dropdown button in Gallery Edit page | Valid `galleryByFragment` configuration. |
| Scrape gallery from URL | Valid `galleryByURL` configuration with matching URL. |

URL-based scraping accepts multiple scrape configurations, and each configuration requires a `url` field. stash iterates through these configurations, attempting to match the entered URL against the `url` fields in the configuration. It executes the first scraping configuration where the entered URL contains the value of the `url` field. 

    
## Actions

### Script

Executes a script to perform the scrape. The `script` field is required for this action and accepts a list of string arguments. For example:

```yaml
action: script
script:
  - python
  - iafdScrape.py
  - query
```

If the script specifies the python executable, Stash will find the correct python executable for your system, either `python` or `python3`. So for example. this configuration could execute `python iafdScrape.py query` or `python3 iafdScrape.py query`.
`python3` will be looked for first and if it's not found, we'll check for `python`. In the case neither are found, you will get an error.

Stash sends data to the script process's `stdin` stream and expects the output to be streamed to the `stdout` stream. Any errors and progress messages should be output to `stderr`.

The script is sent input and expects output based on the scraping type, as detailed in the following table:

| Scrape type | Input | Output |
|-------------|-------|--------|
| `performerByName` | `{"name": "<performer query string>"}` | Array of JSON-encoded performer fragments (including at least `name`) |
| `performerByFragment` | JSON-encoded performer fragment | JSON-encoded performer fragment |
| `performerByURL` | `{"url": "<url>"}` | JSON-encoded performer fragment |
| `sceneByName` | `{"name": "<scene query string>"}` | Array of JSON-encoded scene fragments |
| `sceneByQueryFragment`, `sceneByFragment` | JSON-encoded scene fragment | JSON-encoded scene fragment |
| `sceneByURL` | `{"url": "<url>"}` | JSON-encoded scene fragment |
| `movieByURL` | `{"url": "<url>"}` | JSON-encoded movie fragment |
| `galleryByFragment` | JSON-encoded gallery fragment | JSON-encoded gallery fragment |
| `galleryByURL` | `{"url": "<url>"}` | JSON-encoded gallery fragment |

For `performerByName`, only `name` is required in the returned performer fragments. One entire object is sent back to `performerByFragment` to scrape a specific performer, so the other fields may be included to assist in scraping a performer. For example, the `url` field may be filled in for the specific performer page, then `performerByFragment` can extract by using its value.
  
Python example of a performer Scraper:
  
```python
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

This action scrapes a web page using an xpath configuration to parse. This action is **not valid** for `performerByFragment`.

This action requires that the top-level `xPathScrapers` configuration is populated. The `scraper` field is required and must match the name of a scraper name configured in `xPathScrapers`. For example:

```yaml
sceneByURL:
- action: scrapeXPath
  url: 
    - pornhub.com/view_video.php
  scraper: sceneScraper
```

The above configuration requires that `sceneScraper` exists in the `xPathScrapers` configuration.

XPath scraping configurations specify the mapping between object fields and an xpath selector. The xpath scraper scrapes the applicable URL and uses xpath to populate the object fields.

### scrapeJson

This action works in the same way as `scrapeXPath`, but uses a mapped json configuration to parse. It uses the top-level `jsonScrapers` configuration. This action is **not valid** for `performerByFragment`.

JSON scraping configurations specify the mapping between object fields and a GJSON selector. The JSON scraper scrapes the applicable URL and uses [GJSON](https://github.com/tidwall/gjson/blob/master/SYNTAX.md) to parse the returned JSON object and populate the object fields.


### scrapeXPath and scrapeJson use with `performerByName`

For `performerByName`, the `queryURL` field must be present also. This field is used to perform a search query URL for performer names. The placeholder string sequence `{}` is replaced with the performer name search string. For the subsequent performer scrape to work, the `URL` field must be filled in with the URL of the performer page that matches a URL given in a `performerByURL` scraping configuration. For example:

```yaml
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

### scrapeXPath and scrapeJson use with `sceneByFragment` and `sceneByQueryFragment`

For `sceneByFragment` and `sceneByQueryFragment`, the `queryURL` field must also be present. This field is used to build a query URL for scenes. For `sceneByFragment`, the `queryURL` field supports the following placeholder fields:

* `{checksum}` - the MD5 checksum of the scene
* `{oshash}` - the oshash of the scene
* `{filename}` - the base filename of the scene
* `{title}` - the title of the scene
* `{url}` - the url of the scene

These placeholder field values may be manipulated with regex replacements by adding a `queryURLReplace` section, containing a map of placeholder field to regex configuration which uses the same format as the `replace` post-process action covered below.

For example:

```yaml
sceneByFragment:
  action: scrapeJson
  scraper: sceneQueryScraper
  queryURL: https://metadataapi.net/api/scenes?parse={filename}&limit=1
  queryURLReplace:
    filename:
      - regex: <some regex>
        with: <replacement>
```

The above configuration would scrape from the value of `queryURL`, replacing `{filename}` with the base filename of the scene, after it has been manipulated by the regex replacements.

### scrapeXPath and scrapeJson use with `<scene|performer|gallery|movie>ByURL`

For `sceneByURL`, `performerByURL`, `galleryByURL` the `queryURL` can also be present if we want to use `queryURLReplace`. The functionality is the same as `sceneByFragment`, the only placeholder field available though is the `url`:
* `{url}` - the url of the scene/performer/gallery

```yaml
sceneByURL:
  - action: scrapeJson
    url:
      - metartnetwork.com
    scraper: sceneScraper
    queryURL: "{url}"
    queryURLReplace:
      url:
        - regex: '^(?:.+\.)?([^.]+)\.com/.+movie/(\d+)/(\w+)/?$'
          with: https://www.$1.com/api/movie?name=$3&date=$2
```

### Stash

A different stash server can be configured as a scraping source. This action applies only to `performerByName`, `performerByFragment`, and `sceneByFragment` types. This action requires that the top-level `stashServer` field is configured.

`stashServer` contains a single `url` field for the remote stash server. The username and password can be embedded in this string using `username:password@host`.

An example stash scrape configuration is below:

```yaml
name: stash
performerByName:
  action: stash
performerByFragment:
  action: stash
sceneByFragment:
  action: stash
stashServer:
  url: http://stashserver.com:9999
```
  
## Xpath and JSON scrapers configuration

The top-level `xPathScrapers` field contains xpath scraping configurations, freely named. These are referenced in the `scraper` field for `scrapeXPath` scrapers. 

Likewise, the top-level `jsonScrapers` field contains json scraping configurations.

Collectively, these configurations are known as mapped scraping configurations. 

A mapped scraping configuration may contain a `common` field, and must contain `performer`, `scene`, `movie` or `gallery` depending on the scraping type it is configured for. 

Within the `performer`/`scene`/`movie`/`gallery` field are key/value pairs corresponding to the [golang fields](/help/ScraperDevelopment.md#object-fields) on the performer/scene object. These fields are case-sensitive. 

The values of these may be either a simple selector value, which tells the system where to get the value of the field from, or a more advanced configuration (see below). For example, for an xpath configuration:

```yaml
performer:
  Name: //h1[@itemprop="name"]
```

This will set the `Name` attribute of the returned performer to the text content of the element that matches `<h1 itemprop="name">...`.

For a json configuration:

```yaml
performer:
  Name: data.name
```

The value may also be a sub-object. If it is a sub-object, then the selector must be set to the `selector` key of the sub-object. For example, using the same xpath as above:

```yaml
performer:
  Name: 
    selector: //h1[@itemprop="name"]
    postProcess:
      # post-processing config values
```

### Fixed attribute values

Alternatively, an attribute value may be set to a fixed value, rather than scraping it from the webpage. This can be done by replacing `selector` with `fixed`. For example:

```yaml
performer:
  Gender: 
    fixed: Female
```

### Common fragments

The `common` field is used to configure selector fragments that can be referenced in the selector strings. These are key-value pairs where the key is the string to reference the fragment, and the value is the string that the fragment will be replaced with. For example:

```yaml
common:
  $infoPiece: //div[@class="infoPiece"]/span
performer:
  Measurements: $infoPiece[text() = 'Measurements:']/../span[@class="smallInfo"]
```

The `Measurements` xpath string will replace `$infoPiece` with `//div[@class="infoPiece"]/span`, resulting in: `//div[@class="infoPiece"]/span[text() = 'Measurements:']/../span[@class="smallInfo"]`.

> **⚠️ Note:** Recursive common fragments are **not** supported.  
Referencing a common fragment within another common fragment will cause an error. For example:
```yaml
common:
  $info: //div[@class="info"]
  # Referencing $info in $models will cause an error
  $models: $info/a[@class="model"]
scene:
  Title: $info/h1
  Performers:
    Name: $models
    URL: $models/@href
```

### Post-processing options

Post-processing operations are contained in the `postProcess` key. Post-processing operations are performed in the order they are specified. The following post-processing operations are available:
* `javascript`: accepts a javascript code block, that must return a string value. The input string is declared in the `value` variable. If an error occurs while compiling or running the script, then the original value is returned.
Example:
```yaml
performer:
  Name:
    selector: //div[@class="example element"]
    postProcess:
      - javascript: |
          // capitalise the first letter
          if (value && value.length) {
            return value[0].toUpperCase() + value.substring(1)
          }
```
Note that the `otto` javascript engine is missing a few built-in methods and may not be consistent with other modern javascript implementations.
* `feetToCm`: converts a string containing feet and inches numbers into centimeters. Looks for up to two separate integers and interprets the first as the number of feet, and the second as the number of inches. The numbers can be separated by any non-numeric character including the `.` character. It does not handle decimal numbers. For example `6.3` and `6ft3.3` would both be interpreted as 6 feet, 3 inches before converting into centimeters.
* `lbToKg`: converts a string containing lbs to kg.
* `map`: contains a map of input values to output values. Where a value matches one of the input values, it is replaced with the matching output value. If no value is matched, then value is unmodified.

Example:
```yaml
performer:
  Gender:
    selector: //div[@class="example element"]
    postProcess:
      - map:
          F: Female
          M: Male
  Height:
    selector: //span[@id="height"]
    postProcess:
      - feetToCm: true
  Weight:
    selector: //span[@id="weight"]
    postProcess:
      - lbToKg: true
```
Gets the contents of the selected div element, and sets the returned value to `Female` if the scraped value is `F`; `Male` if the scraped value is `M`.
Height and weight are extracted from the selected spans and converted to `cm` and `kg`.

* `parseDate`: if present, the value is the date format using go's reference date (2006-01-02). For example, if an example date was `14-Mar-2003`, then the date format would be `02-Jan-2006`. See the [time.Parse documentation](https://golang.org/pkg/time/#Parse) for details. When present, the scraper will convert the input string into a date, then convert it to the string format used by stash (`YYYY-MM-DD`). Strings "Today", "Yesterday" are matched (case insensitive) and converted by the scraper so you don't need to edit/replace them. 
Unix timestamps (example: 1660169451) can also be parsed by selecting `unix` as the date format.
Example:
```yaml
Date:
  selector: //div[@class="value epoch"]/text()
  postProcess:
    - parseDate: unix
```

* `subtractDays`: if set to `true` it subtracts the value in days from the current date and returns the resulting date in stash's date format.
Example:
```yaml
Date:
  selector: //strong[contains(text(),"Added:")]/following-sibling::text()
  postProcess:
    - replace:
        - regex: (\d+)\sdays\sago.+
          with: $1
    - subtractDays: true
```

* `replace`: contains an array of sub-objects. Each sub-object must have a `regex` and `with` field. The `regex` field is the regex pattern to replace, and `with` is the string to replace it with. `$` is used to reference capture groups - `$1` is the first capture group, `$2` the second and so on. Replacements are performed in order of the array.

Example:
```yaml
CareerLength: 
  selector: $infoPiece[text() = 'Career Start and End:']/../span[@class="smallInfo"]
    postProcess:
      - replace:
          - regex: \s+to\s+
            with: "-"
```
Replaces `2001 to 2003` with `2001-2003`.

* `subScraper`: if present, the sub-scraper will be executed after all other post-processes are complete and before parseDate. It then takes the value and performs an http request, using the value as the URL. Within the `subScraper` config is a nested scraping configuration. This allows you to traverse to other webpages to get the attribute value you are after. For more info and examples have a look at [#370](https://github.com/stashapp/stash/pull/370), [#606](https://github.com/stashapp/stash/pull/606)

Additionally, there are a number of fixed post-processing fields that are specified at the attribute level (not in `postProcess`) that are performed after the `postProcess` operations:
* `concat`: if an xpath matches multiple elements, and `concat` is present, then all of the elements will be concatenated together
* `split`: the inverse of `concat`. Splits a string to more elements using the separator given. For more info and examples have a look at PR [#579](https://github.com/stashapp/stash/pull/579)

Example:
```yaml
Tags:
  Name:
    selector: //span[@class="list_attributes"]
    split: ","
```
Splits a comma separated list of tags located in the span and returns the tags.


For backwards compatibility, `replace`, `subscraper` and `parseDate` are also allowed as keys for the attribute.

Post-processing on attribute post-process is done in the following order: `concat`, `replace`, `subscraper`, `parseDate` and then `split`.

### XPath resources:

- Test XPaths in Firefox: https://addons.mozilla.org/en-US/firefox/addon/try-xpath/
- XPath cheatsheet: https://devhints.io/xpath

### GJSON resources:

- GJSON Path Syntax: https://github.com/tidwall/gjson/blob/master/SYNTAX.md

### Debugging support
To print the received html/json from a scraper request to the log file, add the following to your scraper yml file:
```yaml
debug:
  printHTML: true
```

### CDP support

Some websites deliver content that cannot be scraped using the raw html file alone. These websites use javascript to dynamically load the content. As such, direct xpath scraping will not work on these websites. There is an option to use Chrome DevTools Protocol to load the webpage using an instance of Chrome, then scrape the result.

Chrome CDP support can be enabled for a specific scraping configuration by adding the following to the root of the yml configuration:
```yaml
driver:
  useCDP: true
```

Optionally, you can add a `sleep` value under the `driver` section. This specifies the amount of time (in seconds) that the scraper should wait after loading the website to perform the scrape. This is needed as some sites need more time for loading scripts to finish. If unset, this value defaults to 2 seconds.

When `useCDP` is set to true, stash will execute or connect to an instance of Chrome. The behavior is dictated by the `Chrome CDP path` setting in the user configuration. If left empty, stash will attempt to find the Chrome executable in the path environment, and will fail if it cannot find one. 

`Chrome CDP path` can be set to a path to the chrome executable, or an http(s) address to remote chrome instance (for example: `http://localhost:9222/json/version`). As remote instance a docker container can also be used with the `chromedp/headless-shell` image being highly recommended.

### CDP Click support

When using CDP you can use  the `clicks` part of the `driver` section to do Mouse Clicks on elements you need to collapse or toggle. Each click element has an `xpath` value that holds the XPath for the button/element you need to click and an optional `sleep` value that is the time in seconds to wait for after clicking.
If the `sleep` value is not set it defaults to `2` seconds.

A demo scraper using `clicks` follows.

```yaml
name: clickDemo # demo only for a single URL
sceneByURL:
  - action: scrapeXPath
    url:
      - https://getbootstrap.com/docs/4.3/components/collapse/
    scraper: sceneScraper

xPathScrapers:
  sceneScraper:
    scene:
      Title: //head/title
      Details: # shows the id/s of the the visible div/s for the Multiple targets example of the page
        selector: //div[@class="bd-example"]//div[@class="multi-collapse collapse show"]/@id
        concat: "\n\n"

driver:
  useCDP: true
  sleep: 1
  clicks: # demo usage toggle on off multiple times
    - xpath: //a[@href="#multiCollapseExample1"] # toggle on first element
    - xpath: //button[@data-target="#multiCollapseExample2"] # toggle on second element
      sleep: 4
    - xpath: //a[@href="#multiCollapseExample1"] # toggle off fist element
      sleep: 1
    - xpath: //button[@data-target="#multiCollapseExample2"] # toggle off second element
    - xpath: //button[@data-target="#multiCollapseExample2"] # toggle on second element
```

> **⚠️ Note:** each `click` adds an extra delay of `clicks sleep` seconds, so the above adds `2+4+1+2+2=11` seconds to the loading time of the page.

### Cookie support

In some websites the use of cookies is needed to bypass a welcoming message or some other kind of protection. Stash supports the setting of cookies for the direct xpath scraper and the CDP based one. Due to implementation issues the usage varies a bit.

To use the cookie functionality a `cookies` sub section needs to be added to the `driver` section.
Each cookie element can consist of a `CookieURL` and a number of `Cookies`.

* `CookieURL` is only needed if you are using the direct / native scraper method. It is the request url that we expect from the site we scrape. It must be in the same domain as the cookies we try to set otherwise all cookies in the same group will fail to set. If the `CookieURL` is not a valid URL then again the cookies of that group will fail.

* `Cookies` are the actual cookies we set. When using CDP that's the only part required. They have  `Name`, `Value`, `Domain`, `Path` values.

In the following example we use cookies for a site using the direct / native xpath scraper. We expect requests to come from `https://www.example.com` and `https://api.somewhere.com` that look for a `_warning` and a `_warn` cookie. A `_test2` cookie is also set just as a demo.

```yaml
driver:
  cookies:
    - CookieURL: "https://www.example.com"
      Cookies:
        - Name: "_warning"
          Domain: ".example.com"
          Value: "true"
          Path: "/"
        - Name: "_test2"
          Value: "123412"
          Domain: ".example.com"
          Path: "/"
    - CookieURL: "https://api.somewhere.com"
      Cookies:
        - Name: "_warn"
          Value: "123"
          Domain: ".somewhere.com"
```

The same functionality when using CDP would look like this:

```yaml
driver:
  useCDP: true
  cookies:
    - Cookies:
        - Name: "_warning"
          Domain: ".example.com"
          Value: "true"
          Path: "/"
        - Name: "_test2"
          Value: "123412"
          Domain: ".example.com"
          Path: "/"
    - Cookies:
        - Name: "_warn"
          Value: "123"
          Domain: ".somewhere.com"
```

For some sites, the value of the cookie itself doesn't actually matter. In these cases, we can use the `ValueRandom`
property instead of `Value`. Unlike `Value`, `ValueRandom` requires an integer value greater than `0` where the value
indicates how long the cookie string should be.

In the following example, we will adapt the previous cookies to use `ValueRandom` instead. We set the `_test2` cookie
to randomly generate a value with a length of 6 characters and the `_warn` cookie to a length of 3.

```yaml
driver:
  cookies:
    - CookieURL: "https://www.example.com"
      Cookies:
        - Name: "_warning"
          Domain: ".example.com"
          Value: "true"
          Path: "/"
        - Name: "_test2"
          ValueRandom: 6
          Domain: ".example.com"
          Path: "/"
    - CookieURL: "https://api.somewhere.com"
      Cookies:
        - Name: "_warn"
          ValueRandom: 3
          Domain: ".somewhere.com"
```

When developing a scraper you can have a look at the cookies set by a site by adding

* a `CookieURL` if you use the direct xpath scraper

* a `Domain` if you use the CDP scraper

and having a look at the log / console in debug mode.

### Headers

Sending request headers is possible when using a scraper.
Headers can be set in the `driver` section and are supported for plain, CDP enabled and JSON scrapers.
They consist of a Key and a Value. If the the Key is empty or not defined then the header is ignored.

```yaml
driver:
  headers:
    - Key: User-Agent
      Value: My Stash Scraper
    - Key: Authorization
      Value: Bearer ds3sdfcFdfY17p4qBkTVF03zscUU2glSjWF17bZyoe8
```

* headers are set after stash's `User-Agent` configuration option is applied.
This means setting a `User-Agent` header from the scraper overrides the one in the configuration settings.

### XPath scraper example

A performer and scene xpath scraper is shown as an example below:

```yaml
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
        postProcess:
          - replace: 
              - regex: .*\((\d+) cm\)
                with: $1
      Ethnicity: $infoPiece[text() = 'Ethnicity:']/../span[@class="smallInfo"]
      FakeTits: $infoPiece[text() = 'Fake Boobs:']/../span[@class="smallInfo"]
      Piercings: $infoPiece[text() = 'Piercings:']/../span[@class="smallInfo"]
      Tattoos: $infoPiece[text() = 'Tattoos:']/../span[@class="smallInfo"]
      CareerLength: 
        selector: $infoPiece[text() = 'Career Start and End:']/../span[@class="smallInfo"]
        postProcess:
          - replace:
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

### JSON scraper example

A performer and scene scraper for ThePornDB is shown below:

```yaml
name: ThePornDB
performerByName:
  action: scrapeJson
  queryURL: https://api.metadataapi.net/performers?q={}
  scraper: performerSearch
performerByURL:
  - action: scrapeJson
    url:
      - https://api.metadataapi.net/performers/
    scraper: performerScraper
sceneByURL:
  - action: scrapeJson
    url:
      - https://api.metadataapi.net/scenes/
    scraper: sceneScraper
sceneByFragment:
  action: scrapeJson
  queryURL: https://api.metadataapi.net/scenes?parse={filename}&hash={oshash}&limit=1
  scraper: sceneQueryScraper
  queryURLReplace:
    filename:
      - regex: "[^a-zA-Z\\d\\-._~]" # clean filename so that it can construct a valid url
        with: "." # "%20"
      - regex: HEVC
        with:
      - regex: x265
        with:
      - regex: \.+
        with: "."
jsonScrapers:
  performerSearch:
    performer:
      Name: data.#.name
      URL:
        selector: data.#.id
        postProcess:
          - replace:
              - regex: ^
                with: https://api.metadataapi.net/performers/

  performerScraper:
    common:
      $extras: data.extras
    performer:
      Name: data.name
      Gender: $extras.gender
      Birthdate: $extras.birthday
      Ethnicity: $extras.ethnicity
      Height:
        selector: $extras.height
        postProcess:
          - replace:
              - regex: cm
                with:
      Measurements: $extras.measurements
      Tattoos: $extras.tattoos
      Piercings: $extras.piercings
      Aliases: data.aliases
      Image: data.image

  sceneScraper:
    common:
      $performers: data.performers
    scene:
      Title: data.title
      Details: data.description
      Date: data.date
      URL: data.url
      Image: data.background.small
      Performers:
        Name: data.performers.#.name
      Studio:
        Name: data.site.name
      Tags:
        Name: data.tags.#.tag

  sceneQueryScraper:
    common:
      $data: data.0
      $performers: data.0.performers
    scene:
      Title: $data.title
      Details: $data.description
      Date: $data.date
      URL: $data.url
      Image: $data.background.small
      Performers:
        Name: $data.performers.#.name
      Studio:
        Name: $data.site.name
      Tags:
        Name: $data.tags.#.tag
driver:
  headers:
    - Key: User-Agent
      Value: Stash JSON Scraper
    - Key: Authorization
      Value: Bearer lPdwFdfY17p4qBkTVF03zscUU2glSjdf17bZyoe  # use an actual API Key here
# Last Updated April 7, 2021
```

## Object fields
### Performer

```
Name
Gender
URL
Twitter
Instagram
Birthdate
DeathDate
Ethnicity
Country
HairColor
EyeColor
Height
Weight
Measurements
FakeTits
CareerLength
Tattoos
Piercings
Aliases
Tags (see Tag fields)
Image
Details
```

*Note:*  - `Gender` must be one of `male`, `female`, `transgender_male`, `transgender_female`, `intersex`, `non_binary` (case insensitive).

### Scene
```
Title
Details
Code
Director
URL
Date
Image
Studio (see Studio Fields)
Movies (see Movie Fields)
Tags (see Tag fields)
Performers (list of Performer fields)
```
### Studio
```
Name
URL
```

### Tag
```
Name
```

### Movie
```
Name
Aliases
Duration
Date
Rating
Director
Studio
Synopsis
URL
FrontImage
BackImage
```

### Gallery
```
Title
Details
URL
Date
Rating
Studio (see Studio Fields)
Tags (see Tag fields)
Performers (list of Performer fields)
```
