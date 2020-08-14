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

## Hashing algorithms

Stash identifies video files by calculating a hash of the file. There are two algorithms available for hashing: `oshash` and `MD5`. `MD5` requires reading the entire file, and can therefore be slow, particularly when reading files over a network. `oshash` (which uses OpenSubtitle's hashing algorithm) only reads 64k from each end of the file.

The hash is used to name the generated files such as preview images and videos, and sprite images.

By default, new systems have MD5 calculation disabled for optimal performance. Existing systems that are upgraded will have the oshash populated for each scene on the next scan. 

### Changing the hashing algorithm

To change the file naming hash to oshash, all scenes must have their oshash values populated. oshash population is done automatically when scanning. 

To change the file naming hash to `MD5`, the MD5 must be populated for all scenes. To do this, `Calculate MD5` for videos must be enabled and the library must be rescanned.

MD5 calculation may only be disabled if the file naming hash is set to `oshash`.

After changing the file naming hash, any existing generated files will now be named incorrectly. This means that stash will not find them and may regenerate them if the `Generate task` is used. To remedy this, run the `Rename generated files` task, which will rename existing generated files to their correct names.

#### Step-by-step instructions to migrate to oshash for existing users

These instructions are for existing users whose systems will be defaulted to use and calculate MD5 checksums. Once completed, MD5 checksums will no longer be calculated when scanning, and oshash will be used for generated file naming. Existing calculated MD5 checksums will remain on scenes, but checksums will not be calculated for new scenes.

1. Scan the library (to populate oshash for all existing scenes).
2. In Settings -> Configuration page, untick `Calculate MD5` and select `oshash` as file naming hash. Save the configuration.
3. In Settings -> Tasks page, click on the `Rename generated files` migration button.

## Scraping 

### User Agent string

Some websites require a legitimate User-Agent string when receiving requests, or they will be rejected. If entered, this string will be applied as the `User-Agent` header value in http scrape requests.

### Chrome CDP path

Some scrapers require a Chrome instance to function correctly. If left empty, stash will attempt to find the Chrome executable in the path environment, and will fail if it cannot find one. 

`Chrome CDP path` can be set to a path to the chrome executable, or an http(s) address to remote chrome instance (for example: `http://localhost:9222/json/version`).

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

