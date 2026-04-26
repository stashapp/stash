host=ghcr.io
user=stashapp
repo=compiler
version=13

VERSION_IMAGE = ${host}/${user}/${repo}:${version}
LATEST_IMAGE = ${host}/${user}/${repo}:latest

latest:
	docker build -t ${LATEST_IMAGE} .

build:
	docker build -t ${VERSION_IMAGE} -t ${LATEST_IMAGE} .

build-no-cache:
	docker build --no-cache -t ${VERSION_IMAGE} -t ${LATEST_IMAGE} .

# requires docker login ghcr.io
# echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
push:
	docker push ${VERSION_IMAGE}
	docker push ${LATEST_IMAGE}