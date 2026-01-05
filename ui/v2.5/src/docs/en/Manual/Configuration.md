# Configuration

## Library

This section enables you to add or remove directories that will be discoverable by Stash. The directories you add will be utilized for scanning new files and for updating their locations in Stash database.

You can configure these directories to apply specifically to:

- **Videos**
- **Images**
- **Both**

> **⚠️ Note:** Don't forget to click `Save` after updating these directories!

## Excluded patterns

Given a valid [regex](https://github.com/google/re2/wiki/Syntax), files that match even partially are excluded during the Scan process and are not entered in the database. Also during the Clean task if these files exist in the DB they are removed from it and their generated files get deleted.  
Prior to matching both the filenames and patterns are converted to lower case so the match is case insensitive.

Regex patterns can be added in the config file or from the UI.  
If you add manually to the config file a restart is needed while from the UI you just need to click the Save button.  
When added through the config file directly special care must be given to double escape the `\` character.

There are 2 separate exclusion settings. One is for videos, another is for images/galleries.

Some examples:

- `"sample\.mp4$"` will exclude all files ending in `sample.mp4`. 
- `"/\.[[:word:]]+/"` will exclude all hidden directories like `/.directoryname/`.
- `"c:\\stash\\videos\\exclude"` will exclude specific Windows directory `c:\stash\videos\exclude`.
- `"^/stash/videos/exclude/"` will exclude all directories that match `/stash/videos/exclude/` pattern.
- `"\\\\stash\\network\\share\\excl\\"` will exclude specific Windows network path `\\stash\network\share\excl\`.

> **Note:** If a directory is excluded for images and videos, then the directory will be excluded from scans completely.

_There is a useful [regex101](https://regex101.com/) site that can help test and experiment with regexps._

## Gallery creation from folders

In the Library section you can find an option to create a gallery from each folder containing images. This will be applied on all libraries when activated, including the base folder of a library. 

If you wish to apply this on a per folder basis, you can create a file called **.nogallery** or **.forcegallery** in a folder that should act different than this global setting.

This will either exclude the folder from becoming a gallery even if the setting is set, or create a gallery from the folder even if the setting is not set. 

The file will only be recognized if written in lower case letters.

Files with a dot in front are handled as hidden in the Linux OS and Mac OS, so you will not see those files after creation on your system without setting your file manager accordingly.

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
2. In Settings -> System page, untick `Calculate MD5` and select `oshash` as file naming hash. Save the configuration.
3. In Settings -> Tasks page, click on the `Rename generated files` migration button.


## Parallel scan/generation

#### Number of parallel task for scan/generation

This setting controls how many sub-tasks will be run in parallel during scanning and generation tasks. (See Tasks)

Auto-detection can be enabled by setting this to zero. This will calculate the number of parallel tasks to be logical cores/4 + 1.

This setting can be used to increase/decrease overall CPU utilisation in two scenarios:

1. High performance 4+ core cpus.
2. Media files stored on remote/cloud filesystem.

Note: If this is set too high it will decrease overall performance and causes failures (out of memory).

## Hardware accelerated live transcoding

Hardware accelerated live transcoding can be enabled by setting the `FFmpeg hardware encoding` setting. Stash outputs the supported hardware encoders to the log file on startup at the Info log level. If a given hardware encoder is not supported, it's error message is logged to the Debug log level for debugging purposes.

## HLS/DASH streaming

To stream using HLS (such as on Apple devices) or DASH, the Cache path must be set. This directory is used to store temporary files during the live-transcoding process. The Cache path can be set in the System settings page. 

## ffmpeg arguments

Additional arguments can be injected into ffmpeg when generating previews and sprites, and when live-transcoding videos. 

The ffmpeg arguments configuration is split into `Input` and `Output` arguments. Input arguments are injected before the input file argument, and output arguments are injected before the output file argument.

Arguments are accepted as a list of strings. Each string is a separate argument. For example, a single argument of `-foo bar` would be treated as a single argument `"-foo bar"`. The correct way to pass this argument would be to split it into two separate arguments: `"-foo", "bar"`.

## Scraping

### User Agent string

Some websites require a legitimate User-Agent string when receiving requests, or they will be rejected. If entered, this string will be applied as the `User-Agent` header value in http scrape requests.

### Chrome CDP path

Some scrapers require a Chrome instance to function correctly. If left empty, stash will attempt to find the Chrome executable in the path environment, and will fail if it cannot find one.

`Chrome CDP path` can be set to a path to the chrome executable, or an http(s) address to remote chrome instance (for example: `http://localhost:9222/json/version`).

> **Important**: As of Chrome 136 you need to specify `--user-data-dir` alongside `--remote-debugging-port`. Read more on their [official post](https://developer.chrome.com/blog/remote-debugging-port). 

## Authentication

By default, stash is not configured with any sort of password protection. To enable password protection, both `Username` and `Password` must be populated. Note that when entering a new username and password where none was set previously, the system will immediately request these credentials to log you in.

## API key

If password protection is enabled, you may also generate an API key. An API key is used by external systems to access your stash system without needing to login first.

External systems using the API key must set the `ApiKey` header value to the configured API key in order to bypass the login requirement.

### Logging out

The logout button is situated in the upper-right part of the screen when you are logged in.

### Recovering from a forgotten username or password

Stash saves login credentials in the config.yml file. You must reset both login and password if you have forgotten your password by doing the following:

* Close your Stash process
* Open the `config.yml` file found in your Stash directory with a text editor
* Delete the `username` and `password` lines from the file and save

Stash authentication should now be reset with no authentication credentials.

## Advanced configuration options

These options are typically not exposed in the UI and must be changed manually in the `config.yml` file.

| Field | Remarks |
|-------|---------|
| `custom_served_folders` | A map of URLs to file system folders. See below. |
| `custom_ui_location` | The file system folder where the UI files will be served from, instead of using the embedded UI. Empty to disable. Stash must be restarted to take effect. |
| `developer_options.extra_blob_paths` | A list of alternative blob paths. These paths will be read for blob files. Blobs will not be written or deleted from these paths. Intended for developer use only. |
| `max_upload_size` | Maximum file upload size for import files. Defaults to 1GB. |
| `theme_color` | Sets the `theme-color` property in the UI. |
| `gallery_cover_regex` | The regex responsible for selecting images as gallery covers |
| `proxy` | The url of a HTTP(S) proxy to be used when stash makes calls to online services. Example: https://user:password@my.proxy:8080. Note: SOCKS5 proxies are unsupported. |
| `no_proxy` | A list of domains for which the proxy must not be used. Default is all local LAN: localhost,127.0.0.1,192.168.0.0/16,10.0.0.0/8,172.16.0.0/12 |
| `sequential_scanning` | Modifies behaviour of the scanning functionality to generate support files (previews/sprites/phash) at the same time as fingerprinting/screenshotting. Useful when scanning cached remote files. |

The following environment variables are also supported:

| Environment variable | Remarks |
|----------------------|---------|
| `STASH_SQLITE_CACHE_SIZE` | Sets the SQLite cache size. See https://www.sqlite.org/pragma.html#pragma_cache_size. Default is `-2000` which is 2MB. |

### Custom favicon

You can provide a custom favicon by placing a `favicon.ico` file in the configuration directory. The configuration directory is located alongside the `config.yml` file.

When a custom favicon is provided, it will be served instead of the default embedded favicon.

### Custom served folders

Custom served folders are served when the server handles a request with the `/custom` URL prefix. The following is an example configuration:

```
custom_served_folders:
  /: D:\stash\static
  /foo: D:\bar
```

With the above configuration, a request for `/custom/foo/bar.png` would serve `D:\bar\bar.png`. 

The `/` entry matches anything that is not otherwise mapped by the other entries. For example, `/custom/baz/xyz.png` would serve `D:\stash\static\baz\xyz.png`.
