lint:
	@docker run -t --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run -v
