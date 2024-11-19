user=stashapp
repo=compiler
version=10

latest:
	docker build -t ${user}/${repo}:latest .

build:
	docker build -t ${user}/${repo}:${version} -t ${user}/${repo}:latest .

build-no-cache:
	docker build --no-cache -t ${user}/${repo}:${version} -t ${user}/${repo}:latest .

install: build
	docker push ${user}/${repo}:${version}
	docker push ${user}/${repo}:latest
