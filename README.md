# Stash
https://stashapp.cc

[![Build](https://github.com/stashapp/stash/actions/workflows/build.yml/badge.svg?branch=develop&event=push)](https://github.com/stashapp/stash/actions/workflows/build.yml)
[![Docker pulls](https://img.shields.io/docker/pulls/stashapp/stash.svg)](https://hub.docker.com/r/stashapp/Stash 'DockerHub')
[![Go Report Card](https://goreportcard.com/badge/github.com/stashapp/stash)](https://goreportcard.com/report/github.com/stashapp/stash)
[![Discord](https://img.shields.io/discord/559159668438728723.svg?logo=discord)](https://discord.gg/2TsNFKt)

### **Stash is a self-hosted webapp written in Go which organizes and serves your porn.**
![demo image](docs/readme_assets/demo_image.png)

* Stash gathers information about videos in your collection from the internet, and is extensible through the use of community-built plugins for a large number of content producers and sites.
* Stash supports a wide variety of both video and image formats.
* You can tag videos and find them later.
* Stash provides statistics about performers, tags, studios and more.

You can [watch a SFW demo video](https://vimeo.com/545323354) to see it in action.

For further information you can [read the in-app manual](ui/v2.5/src/docs/en).

# Installing Stash

<img src="docs/readme_assets/windows_logo.svg" width="100%" height="75"> Windows | <img src="docs/readme_assets/mac_logo.svg" width="100%" height="75"> MacOS| <img src="docs/readme_assets/linux_logo.svg" width="100%" height="75"> Linux | <img src="docs/readme_assets/docker_logo.svg" width="100%" height="75"> Docker
:---:|:---:|:---:|:---:
[Latest Release](https://github.com/stashapp/stash/releases/latest/download/stash-win.exe) <br /> <sup><sub>[Development Preview](https://github.com/stashapp/stash/releases/download/latest_develop/stash-win.exe)</sub></sup> | [Latest Release (Apple Silicon)](https://github.com/stashapp/stash/releases/latest/download/stash-osx-applesilicon) <br /> <sup><sub>[Development Preview (Apple Silicon)](https://github.com/stashapp/stash/releases/download/latest_develop/stash-osx-applesilicon)</sub></sup> <br>[Latest Release (Intel)](https://github.com/stashapp/stash/releases/latest/download/stash-osx) <br /> <sup><sub>[Development Preview (Intel)](https://github.com/stashapp/stash/releases/download/latest_develop/stash-osx)</sub></sup> | [Latest Release (amd64)](https://github.com/stashapp/stash/releases/latest/download/stash-linux) <br /> <sup><sub>[Development Preview (amd64)](https://github.com/stashapp/stash/releases/download/latest_develop/stash-linux)</sub></sup> <br /> [More Architectures...](https://github.com/stashapp/stash/releases/latest) | [Instructions](docker/production/README.md) <br /> <sup><sub> [Sample docker-compose.yml](docker/production/docker-compose.yml)</sub></sup>

## Installation Notes

#### FFMPEG
Stash requires ffmpeg. If you don't have it installed, Stash will download a copy for you. It is recommended that Linux users install `ffmpeg` from their distro's package manager.
#### Windows Users: Security Prompt
*Note for Windows users:* Running the app might present a security prompt since the binary isn't yet signed.  Bypass this by clicking "more info" and then the "run anyway" button.

# Usage

## Quickstart Guide
### Launching
Once downloaded, Run the executable (double click the exe on windows or run `./stash-osx` / `./stash-linux` from the terminal on macOS / Linux) to get started.

Stash is a web-based application.  Once the application is running, it is availble (by default) from http://localhost:9999.  You can change this in the settings.

### Configuring and Importing Your Content
On first run, Stash will prompt you for some configuration options and a directory to index (you can also skip this step and complete afterward)

### Gathering Content Metadata
After indexing, you can manually view, edit and tag your content as you'd like. 

Stash can pull metadata (performer, description, studio, etc) information from the internet through the use of small snippets of code called  [scrapers](https://github.com/stashapp/stash/tree/develop/ui/v2.5/src/docs/en/Scraping.md) that interpret external database and website information and import it into your program. Stash comes pre-installed with one such scraper.

In order to maximize this functionality, you will need to download Scrapers from the [Community Scrapers Collection](https://github.com/stashapp/CommunityScrapers). (maintained by the Stash community).  Start with a few and add to them as necessary. After downloading, copy the contents of the archive to the /scrapers/ subdirectory of your Stash installation, then restart Stash.  The scraper list will show up in the "Metadata Providers" section of your configuration screen.

The simplest way to tag a large number of files is by using the [Tagger](https://github.com/stashapp/stash/blob/develop/ui/v2.5/src/docs/en/Tagger.md) view which uses filename keywords to help identify the file and pull in scene and performer information from our stash-box database as well as other sources. Note that Stash-Box not comprehensive and you will likely need to use the scrapers to identify some of your media.

# Translation
[![Translate](https://translate.stashapp.cc/widgets/stash/-/stash-desktop-client/svg-badge.svg)](https://translate.stashapp.cc/engage/stash/)
🇧🇷 🇨🇳 🇬🇧 🇫🇮 🇫🇷 🇩🇪 🇮🇹 🇪🇸 🇸🇪 🇹🇼 🇹🇷

Stash is available in 11 languages (so far!) and it could be in your language too. If you want to help us translate Stash into your language, you can make an account at [translate.stashapp.cc](https://translate.stashapp.cc/projects/stash/stash-desktop-client/) to get started contributing new languages or improving existing ones. Thanks!

# Support (FAQ)

Answers to other Frequently Asked Questions can be found [on our Wiki](https://github.com/stashapp/stash/wiki/FAQ)

For issues not addressed there, there are a few options.

* Read the [Wiki](https://github.com/stashapp/stash/wiki)
* Check the in-app documentation, in the top right corner of the app (also available [here](https://github.com/stashapp/stash/tree/develop/ui/v2.5/src/docs/en)
* Join the [Discord server](https://discord.gg/2TsNFKt), where the community can offer support.

# Customization

## Themes and CSS Customization
There is a [directory of community-created themes](https://github.com/stashapp/stash/wiki/Themes) on our Wiki, along with instructions on how to install them.

You can also make Stash interface fit your desired style with [Custom CSS snippets](https://github.com/stashapp/stash/wiki/Custom-CSS-snippets).

# For Developers

Pull requests are welcome! 

See [Development](docs/DEVELOPMENT.md) and [Contributing](docs/CONTRIBUTING.md) for information on working with the codebase, getting a local development setup, and contributing changes.
