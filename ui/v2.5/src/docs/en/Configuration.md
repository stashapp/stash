# Configuration

## Stashes

This section allows you to add and remove directories from your library list. Files in these directories will be included when scanning. Files that are outside of these directories will be removed when running the Clean task.

> **⚠️ Note:** Don't forget to click `Save` after updating these directories!

## Excluded Patterns

Given a valid [regex](https://github.com/google/re2/wiki/Syntax), files that match even partially are excluded during the Scan process and are not entered in the database. Also during the Clean task if these files exist in the DB they are removed from it and their generated files get deleted.
Prior to matching both the filenames and patterns are converted to lower case so the match is case insensitive.

Regex patterns can be added in the config file or from the UI.
If you add manually to the config file a restart is needed while from the UI you just need to click the Save button.
When added through the config file directly special care must be given to double escape the `\` character.

Some examples

For the config file you need the following added
```
exclude: 
- "sample\\.mp4$"
- "/\\.[[:word:]]+/"
- "c:\\\\stash\\\\videos\\\\exclude"
- "^/stash/videos/exclude/"
- "\\\\\\\\stash\\network\\\\share\\\\excl\\\\"
```
* the first excludes all files ending in `sample.mp4` ( `.` needs to be escaped also)
* the second hidden directories `/.directoryname/`
* the third is an example for a windows directory `c:\stash\videos\exclude`
* the fourth the directory `/stash/videos/exclude/`
* and the last a windows network path `\\stash\network\share\excl\`

_a useful [link](https://regex101.com/) to experiment with regexps_

## Scraping User Agent string

Some websites require a legitimate User-Agent string when receiving requests, or they will be rejected. If entered, this string will be applied as the `User-Agent` header value in http scrape requests.

## Authentication

By default, stash is not configured with any sort of password protection. To enable password protection, both `Username` and `Password` must be populated. Note that when entering a new username and password where none was set previously, the system will immediately request these credentials to log you in.

### Logging out

The logout button is situated in the upper-right part of the screen when you are logged in.

### Recovering from a forgotten username or password

Stash saves login credentials in the config.yml file. You must reset both login and password if you have forgotten your password by doing the following:
* Close your Stash process
* Open the `config.yml` file found in your Stash directory with a text editor
* Delete the `login` and `password` lines from the file and save
Stash authentication should now be reset with no authentication credentials.

