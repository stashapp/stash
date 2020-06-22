# Ways to contribute

## Financial

Financial contributions are welcomed and are accepted using [Open Collective](https://opencollective.com/stashapp).

## Development-related

The Stash backend is written in golang with a sqlite database. The UI is written in react. Bug fixes, improvements and new features are welcomed. Please see the [README.md](https://github.com/stashapp/stash/raw/develop/README.md) file for details on how to get started. Assistance can be provided via our [Discord](https://discord.gg/2TsNFKt).

## Documentation

Efforts to improve documentation in stash helps new users and reduces the amount of questions we have to field in Discord. Contributions to documentation are welcomed. While submitting documentation changes via git pull requests is ideal, we will gladly accept submissions via [github issues](https://github.com/stashapp/stash/issues) or on [Discord](https://discord.gg/2TsNFKt).

For those with web page experience, we also welcome contributions to our [website](https://stashapp.cc/) (which as of writing is very undeveloped).

## Testing features, improvements and bug fixes

Testing is currently covered by a very small group, so new testers are welcomed. Being able to build stash locally is ideal, but custom binaries for pull requests are available by navigating to the `continuous-integration/travis-ci/pr` travis check details. 

The link to the custom binary for each platform can be found at the end of the build log, and looks like the following:
```
$ if [ "$TRAVIS_PULL_REQUEST" != "false" ]; then sh ./scripts/upload-pull-request.sh; fi
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 43.1M  100    35  100 43.1M      3  3812k  0:00:11  0:00:11 --:--:-- 5576k
stash-osx uploaded to url: https://transfer.sh/.../stash-osx
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 60.7M  100    39  100 60.7M      3  5391k  0:00:13  0:00:11  0:00:02 7350k
stash-win.exe uploaded to url: https://transfer.sh/.../stash-win.exe
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 44.6M  100    37  100 44.6M      2  3648k  0:00:18  0:00:12  0:00:06 7504k
stash-linux uploaded to url: https://transfer.sh/.../stash-linux

```
The `if` line will need to be expanded to see the details.

## Submitting and contributing to bug reports, improvements and new features

We welcome contributions for future improvements and features, and bug reports help everyone. These can all be found in the [github issues](https://github.com/stashapp/stash/issues).


## Providing support

Offering support for new users on [Discord](https://discord.gg/2TsNFKt) is also welcomed.