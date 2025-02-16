-include .env

START_LOG = @echo "=============== START OF LOG ==============="
END_LOG = @echo "=============== END OF LOG ==============="

.PHONY: env
env:
	$(START_LOG)
	@cp ./.env.tmpl ./.env
	@echo "Environment file created at ./.env.develop"
	$(END_LOG)

.PHONY: infra
infra:
	$(START_LOG)
	@cd third_party/cartesi-coprocessor; docker compose -f docker-compose-devnet.yaml up --build
	$(END_LOG)

.PHONY: local
local:
	$(START_LOG)
	@forge script ./contracts/script/DeployLocal.s.sol --broadcast \
									 --root contracts \
									 --rpc-url $(RPC_URL) \
									 --private-key $(PRIVATE_KEY) \
									 -v
	$(END_LOG)

.PHONY: task_manager
task_manager:
	$(START_LOG)
	@forge script ./contracts/script/DeployTaskManager.s.sol --broadcast \
									 --root contracts \
									 --rpc-url $(RPC_URL) \
									 --private-key $(PRIVATE_KEY) \
									 -v
	$(END_LOG)

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