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
[Latest Release](https://github.com/stashapp/stash/releases/latest/download/stash-win.exe) <br /> <sup><sub>[Development Preview](https://github.com/stashapp/stash/releases/download/latest_develop/stash-win.exe)</sub></sup> | [Latest Release](https://github.com/stashapp/stash/releases/latest/download/Stash-macos.zip) <br /> <sup><sub>[Development Preview](https://github.com/stashapp/stash/releases/download/latest_develop/Stash-macos.zip)</sub></sup> | [Latest Release (amd64)](https://github.com/stashapp/stash/releases/latest/download/stash-linux) <br /> <sup><sub>[Development Preview (amd64)](https://github.com/stashapp/stash/releases/download/latest_develop/stash-linux)</sub></sup> <br /> [More Architectures...](https://github.com/stashapp/stash/releases/latest) | [Instructions](docker/production/README.md) <br /> <sup><sub> [Sample docker-compose.yml](docker/production/docker-compose.yml)</sub></sup>

## Getting Started
Run the executable by double-clicking it.

*Note for Mac users:* Running the app will present a security prompt since the binary isn't yet signed. Bypass this by right (ctrl) clicking on the app and hitting "Open", then "Open" in the next popup.

*Note for Windows users:* Running the app might present a security prompt since the binary isn't yet signed. Bypass this by clicking "more info" and then the "run anyway" button.

#### FFMPEG
Stash requires ffmpeg. If you don't have it installed, Stash will download a copy for you. It is recommended that Linux users install `ffmpeg` from their distro's package manager.

# Usage

## Quickstart Guide
Download and run Stash. It will prompt you for some configuration options and a directory to index (you can also do this step afterward)

**If you'd like to automatically retrieve and organize information about your entire library,**  You will need to download some [scrapers](https://github.com/stashapp/stash/tree/develop/ui/v2.5/src/docs/en/Scraping.md). The Stash community has developed scrapers for many popular data sources which can be downloaded and installed from [this repository](https://github.com/stashapp/CommunityScrapers).

The simplest way to tag a large number of files is by using the [Tagger](https://github.com/stashapp/stash/blob/develop/ui/v2.5/src/docs/en/Tagger.md) which uses filename keywords to help identify the file and pull in scene and performer information from our stash-box database. Note that this data source is not comprehensive and you may need to use the scrapers to identify some of your media.

# Translation
[![Translate](https://translate.stashapp.cc/widgets/stash/-/stash-desktop-client/svg-badge.svg)](https://translate.stashapp.cc/engage/stash/)
🇧🇷 🇨🇳 🇳🇱 🇬🇧 🇫🇮 🇫🇷 🇩🇪 🇮🇹 🇯🇵 🇪🇸 🇸🇪 🇹🇼 🇹🇷

Stash is available in 13 languages (so far!) and it could be in your language too. If you want to help us translate Stash into your language, you can make an account at [translate.stashapp.cc](https://translate.stashapp.cc/projects/stash/stash-desktop-client/) to get started contributing new languages or improving existing ones. Thanks!

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
