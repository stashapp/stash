#phony dependency task that does nothing
#"make executable" does not run if there is a ./executable directory, unless the task has a dependency
phony:

lint:
	docker run --rm -v ${CURDIR}:/app -w /app golangci/golangci-lint:v1.23.8 golangci-lint run -v
