host=ghcr.io
user=stashapp
repo=compiler
version=12

VERSION_IMAGE = ${host}/${user}/${repo}:${version}
LATEST_IMAGE = ${host}/${user}/${repo}:latest

latest:
	docker build -t ${LATEST_IMAGE} .

latest:
	docker build -t ${LATEST_IMAGE} .

build:
	docker build -t ${VERSION_IMAGE} -t ${LATEST_IMAGE} .

build-no-cache:
	docker build --no-cache -t ${VERSION_IMAGE} -t ${LATEST_IMAGE} .