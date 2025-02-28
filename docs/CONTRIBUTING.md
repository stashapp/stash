## Goals and design vision

The goal of stash is to be:
- an application for organising and viewing adult content - currently this is videos and images, in future this will be extended to include audio and text content
  - organising includes scraping of metadata from websites and metadata repositories
- free and open-source
- portable and offline - can be run on a USB stick without needing to install dependencies (with the exception of ffmpeg)
- minimal, but highly extensible. The core feature set should be the minimum required to achieve the primary goal, while being extensible enough to extend via plugins
- easy to learn and use, with minimal technical knowledge required

The core stash system is not intended for:
- managing downloading of content
- managing content on external websites
- publically sharing content

Other requirements:
- support as many video and image formats as possible
- interfaces with external systems (for example stash-box) should be made as generic as possible. 

Design considerations:
- features are easy to add and difficult to remove. Large superfluous features should be scrutinised and avoided where possible (eg DLNA, filename parser). Such features should be considered for third-party plugins instead.

## Technical Debt
Please be sure to consider how heavily your contribution impacts the maintainability of the project long term, sometimes less is more.  We don't want to merge collossal pull requests with hundreds of dependencies by a driveby contributor.

## Contributor Checklist
Please make sure that you've considered the following before you submit your Pull Requests as ready for merging.
* I've run Code linters and [gofmt](https://golang.org/cmd/gofmt/) to make sure that my code is readable.
* I have read through formerly submitted [pull requests](https://github.com/stashapp/stash/pulls) and [git issues](https://github.com/stashapp/stash/issues) to make sure that this contribution is required and isn't a duplicate. Also, so that I can manage to close any git Issues needing closed relating to this feature submission.
* I  commented adequately on my code with the expectation in mind that anyone else should be able to look at this code I've submitted and know exactly what's happening and what the expectations are.

### Legal Agreements
* I acknowledge that if applicable to me, submitting and subsequent acceptance of this Pull Request I, the code contributor of this Pull Request, agree and acknowledge my understanding that the new code license has now been updated to [AGPL](/LICENSE.md). I agree that all code before this Pull Request, which I've previously submitted, is now to be re-licensed under the new license AGPL and no longer the former MIT license.

**In case you were unable to follow any of the above include an explanation as to why not in your Pull Request.**
