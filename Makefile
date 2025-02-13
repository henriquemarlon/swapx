-include .env.develop

START_LOG = @echo "=============== START OF LOG ==============="
END_LOG = @echo "=============== END OF LOG ==============="

.PHONY: env
env: ./.env.develop
	$(START_LOG)
	@cp ./.env.develop.tmpl ./.env.develop
	@touch .cartesi.env
	@echo "Environment file created at ./.env.develop"
	$(END_LOG)

.PHONY: build
build:
	$(START_LOG)
	@docker build \
		-t machine:latest \
		-f ./build/Dockerfile.machine .
	@cartesi build --from-image machine:latest
	$(END_LOG)

.PHONY: dev
dev:
	$(START_LOG)
	@nonodox -- air

.PHONY: test
test:
	@go test -p=1 ./... -coverprofile=./coverage.md -v

.PHONY: coverage
coverage: test
	@go tool cover -html=./coverage.md

.PHONY: state
state:
	@chmod +x ./tools/state.sh
	@./tools/state.sh