.PHONY: cmd
.DEFAULT_GOAL := help

build-image: ## build docker image
	docker build -t ghcr.io/hmasataka/open-match-frontend:latest -f cmd/gamefront/Dockerfile .
	docker build -t ghcr.io/hmasataka/open-match-matchfunction:latest -f cmd/mmf/Dockerfile .
	docker build -t ghcr.io/hmasataka/open-match-director:latest -f cmd/director/Dockerfile .

push-image: ## push docker image
	docker push ghcr.io/hmasataka/open-match-frontend:latest
	docker push ghcr.io/hmasataka/open-match-matchfunction:latest
	docker push ghcr.io/hmasataka/open-match-director:latest

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
