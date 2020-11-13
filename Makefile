export CGO_ENABLED=0

.DEFAULT_GOAL := help
.PHONY: help
help: ## Get help on a command
	@echo '  see: https://github.com/kmrok/image-moderation-poc'
	@echo ''
	@grep -E '^[%/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ./cmd/bin/app ./cmd/main.go

.PHONY: run
run: ## Build docker image and run the app
	@docker build -t image-moderation-poc . > /dev/null
	@$(eval GCLOUD_CREDENTIALS := $(shell cat ./.gcloud/credentials.json))
	@docker run -it --rm  \
		-e 'AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}' \
		-e 'AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}' \
		-e 'GCLOUD_CREDENTIALS=${GCLOUD_CREDENTIALS}' \
		image-moderation-poc
