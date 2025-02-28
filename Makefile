-include .env

START_LOG = @echo "=============== START OF LOG ==============="
END_LOG = @echo "================ END OF LOG ================"

.PHONY: env
env:
	$(START_LOG)
	@cp ./.env.tmpl ./.env
	@echo "Environment file created at ./.env.develop"
	$(END_LOG)

.PHONY: infra
infra:
	$(START_LOG)
	@docker build \
		-f ./third_party/cartesi-coprocessor-operator/Dockerfile \
		-t operator:latest ./third_party/cartesi-coprocessor-operator
	@cd third_party/cartesi-coprocessor; docker compose -f docker-compose-devnet.yaml up --build

.PHONY: clean
clean:
	$(START_LOG)
	@rm -rf ./contracts/out
	@rm -rf ./contracts/cache
	@rm -rf ./contracts/broadcast
	@rm -rf ./third_party/cartesi-coprocessor/operator1-data
	@rm -rf ./third_party/cartesi-coprocessor/env/eigenlayer/anvil/devnet-operators-ready.flag
	@rm -rf ./output.car
	@rm -rf ./output.car.json
	@rm -rf ./output.cid
	@rm -rf ./output.size
	@echo "Cleaned up"
	$(END_LOG)

.PHONY: test
test:
	$(START_LOG)
	@echo "Performing tests..."
	@go test -p=1 ./... -coverprofile=./coverage.md -v
	@forge test --root contracts
	$(END_LOG)

.PHONY: coverage
coverage: test
	$(START_LOG)
	@echo "Reporting coverage..."
	@go tool cover -html=./coverage.md
	$(END_LOG)

.PHONY: fmt
fmt:
	$(START_LOG)
	@echo "Formatting project..."
	@go fmt ./...
	$(END_LOG)

.PHONY: gen
gen:
	$(START_LOG)
	@echo "Generating files..."
	@go generate ./...
	$(END_LOG)

.PHONY: state
state:
	$(START_LOG)
	@chmod +x ./tools/state.sh
	@./tools/state.sh
	$(END_LOG)

.PHONY: slot
slot:
	$(START_LOG)
	@rm -rf storage-layout
	@forge inspect ./src/SwapXHook.sol:SwapXHook storage-layout --root contracts >> storage-layout
	@echo "Storage layout for SwapXHook.sol created at: ./storage-layout"
	$(END_LOG)

.PHONY: v4
v4:
	$(START_LOG)
	@forge script ./contracts/script/V4Deployer.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		-v
	$(END_LOG)

.PHONY: hook
hook:
	$(START_LOG)
	@forge script ./contracts/script/DeployHook.s.sol --broadcast \
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

.PHONY: buy
buy:
	$(START_LOG)
	@forge script ./contracts/script/SendBuyOrder.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		-v
	$(END_LOG)

.PHONY: sell
sell:
	$(START_LOG)
	@forge script ./contracts/script/SendSellOrder.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		-v
	$(END_LOG)