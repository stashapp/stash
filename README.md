# Stash
https://stashapp.cc

[![Build](https://github.com/stashapp/stash/actions/workflows/build.yml/badge.svg?branch=develop&event=push)](https://github.com/stashapp/stash/actions/workflows/build.yml)
[![Docker pulls](https://img.shields.io/docker/pulls/stashapp/stash.svg)](https://hub.docker.com/r/stashapp/stash 'DockerHub')
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
[Latest Release](https://github.com/stashapp/stash/releases/latest/download/stash-win.exe) <br /> <sup><sub>[Development Preview](https://github.com/stashapp/stash/releases/download/latest_develop/stash-win.exe)</sub></sup> | [Latest Release (Apple Silicon)](https://github.com/stashapp/stash/releases/latest/download/stash-macos-applesilicon) <br /> <sup><sub>[Development Preview (Apple Silicon)](https://github.com/stashapp/stash/releases/download/latest_develop/stash-macos-applesilicon)</sub></sup> <br />[Latest Release (Intel)](https://github.com/stashapp/stash/releases/latest/download/stash-macos-intel) <br /> <sup><sub>[Development Preview (Intel)](https://github.com/stashapp/stash/releases/download/latest_develop/stash-macos-intel)</sub></sup> | [Latest Release (amd64)](https://github.com/stashapp/stash/releases/latest/download/stash-linux) <br /> <sup><sub>[Development Preview (amd64)](https://github.com/stashapp/stash/releases/download/latest_develop/stash-linux)</sub></sup> <br /> [More Architectures...](https://github.com/stashapp/stash/releases/latest) | [Instructions](docker/production/README.md) <br /> <sup><sub> [Sample docker-compose.yml](docker/production/docker-compose.yml)</sub></sup>

## First Run
#### Windows Users: Security Prompt
Running the app might present a security prompt since the binary isn't yet signed.  Bypass this by clicking "more info" and then the "run anyway" button.
#### FFMPEG
Stash requires ffmpeg. If you don't have it installed, Stash will download a copy for you. It is recommended that Linux users install `ffmpeg` from their distro's package manager.

# Usage

## Quickstart Guide
Stash is a web-based application. Once the application is running, the interface is available (by default) from http://localhost:9999.

On first run, Stash will prompt you for some configuration options and media directories to index, called "Scanning" in Stash. After scanning, your media will be available for browsing, curating, editing, and tagging.

Stash can pull metadata (performers, tags, descriptions, studios, and more) directly from many sites through the use of [scrapers](https://github.com/stashapp/stash/tree/develop/ui/v2.5/src/docs/en/Scraping.md), which integrate directly into Stash.

Many community-maintained scrapers are available for download at the [Community Scrapers Collection](https://github.com/stashapp/CommunityScrapers). The community also maintains StashDB, a crowd-sourced repository of scene, studio, and performer information, that can automatically identify much of a typical media collection. Inquire in the Discord for details. Identifying an entire collection will typically require a mix of multiple sources. 

<sub>StashDB is the canonical instance of our open source metadata API, [stash-box](https://github.com/stashapp/stash-box).</sub>

# Translation
[![Translate](https://translate.stashapp.cc/widgets/stash/-/stash-desktop-client/svg-badge.svg)](https://translate.stashapp.cc/engage/stash/)
🇧🇷 🇨🇳 🇩🇰 🇳🇱 🇬🇧 🇫🇮 🇫🇷 🇩🇪 🇮🇹 🇯🇵 🇵🇱 🇪🇸 🇸🇪 🇹🇼 🇹🇷

Stash is available in 15 languages (so far!) and it could be in your language too. If you want to help us translate Stash into your language, you can make an account at [translate.stashapp.cc](https://translate.stashapp.cc/projects/stash/stash-desktop-client/) to get started contributing new languages or improving existing ones. Thanks!

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
