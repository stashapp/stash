# Contributing to Stash

## AI Usage Policy

Please see our [AI Usage Policy](/docs/AI_POLICY.md) for guidelines on the use of AI in contributions to this project.

## Issues

Bug reports and feature requests must use descriptive and concise titles and follow the provided templates. Please use the search function to make sure that you are not submitting duplicates, and that a similar report or request has not already been resolved or rejected.

All issues must be written by humans. Fully AI-generated issues will be closed without comment.

## Pull Requests

All pull requests must use descriptive and concise titles and follow the provided templates. In addition, they must follow the the following guidelines:

- You must link to an open issue that pull request addresses (see [GitHub documentation](https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue) on how to do that).
- Pull requests must be focused on a single issue or feature. Large, multi-purpose pull requests will be rejected.
- Large features must be discussed with maintainers before submitting a pull request to ensure it fits with the overall design vision of the project. Failure to do so may result in the pull request being rejected.
- Pull requests must include code tests that sufficiently cover the changes made.
- You must detail the manual testing done and describe the steps taken to sufficiently verify the changes.
- You must be able to explain any line of code and design decision during the review process.
- You may not have more than 3 open pull requests at a time. If you have more than 3 open pull requests, maintainers may ask you to close some of them before they will review any of them.

By submitting a pull request, you agree that you have read and understood and that you are in compliance with the guidelines outlined here, including the [AI Usage Policy](docs/AI_POLICY.md). 

You also agree to license your contribution under the [AGPL](/LICENSE.md) license, and that all of your previous contributions to the project are also licensed under the AGPL.

## Bounties

Pull requests for bounties must be discussed with maintainers before submitting a pull request to ensure it fits with the overall design vision of the project. Failure to do so may result in the pull request being rejected.

## Goals and Design Vision

The goal of Stash is to be:
- An application for organising and viewing NSFW and SFW content - currently this is videos and images, in future this will be extended to include audio and text content
  - Organising includes scraping of metadata from websites and metadata repositories
- Free and open-source
- Portable and offline - can be run on a USB stick without needing to install dependencies (with the exception of FFmpeg)
- Minimal, but highly extensible. The core feature set should be the minimum required to achieve the primary goal, while being extensible enough to extend via plugins
- Easy to learn and use, with minimal technical knowledge required

The core Stash system is not intended for:
- Managing downloading of content
- Managing content on external websites
- Publicly sharing content

Other requirements:
- Support as many video and image formats as possible
- Interfaces with external systems (for example stash-box) should be made as generic as possible. 

Design considerations:
- Features are easy to add and difficult to remove. Large superfluous features should be scrutinised and avoided where possible (e.g. DLNA, filename parser). Such features should be considered for third-party plugins instead.
