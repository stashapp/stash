Modified from https://github.com/bep/dockerfiles/tree/master/ci-goreleaser

When the Dockerfile is changed, the version number should be incremented in the Makefile and the new version tag should be pushed to Docker Hub. The GitHub workflow files also need to be updated to pull the correct image tag.
