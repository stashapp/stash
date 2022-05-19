#!/bin/bash
#
# This was adapted from https://github.com/dgraph-io/dgraph/blob/master/wiki/scripts/build.sh
#

set -e

GREEN='\033[32;1m'
RESET='\033[0m'
HOST=https://gqlgen.com

VERSIONS_ARRAY=(
    'v0.16.0'
    'v0.15.1'
    'origin/master'
    'v0.15.0'
    'v0.14.0'
    'v0.13.0'
    'v0.12.2'
    'v0.11.3'
    'v0.10.2'
    'v0.9.3'
    'v0.8.3'
)

joinVersions() {
	versions=$(printf ",%s" "${VERSIONS_ARRAY[@]}" | sed 's/origin\/master/master/')
	echo "${versions:1}"
}

function version { echo "$@" | gawk -F. '{ printf("%03d%03d%03d\n", $1,$2,$3); }'; }

rebuild() {
	VERSION_STRING=$(joinVersions)
	export CURRENT_VERSION=${1}
	if [[ $CURRENT_VERSION == 'origin/master' ]] ; then
	    CURRENT_VERSION="master"
    fi

	export VERSIONS=${VERSION_STRING}

    hugo --quiet --destination="public/$CURRENT_VERSION" --baseURL="$HOST/$CURRENT_VERSION/"

    if [[ $1 == "${VERSIONS_ARRAY[0]}" ]]; then
        hugo --quiet --destination=public/ --baseURL="$HOST/"
    fi
}


currentBranch=$(git rev-parse --abbrev-ref HEAD)

if ! git remote  | grep -q origin ; then
    git remote add origin https://github.com/99designs/gqlgen
fi
git fetch origin --tags

for version in "${VERSIONS_ARRAY[@]}" ; do
    echo -e "$(date) $GREEN Updating docs for $version.$RESET"
    rm -rf content
    git checkout $version -- content
    rebuild "$version"
done

rm -rf content
git checkout "$currentBranch" -- content

