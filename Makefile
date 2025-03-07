-include .env

.PHONY: env
env:
	@cp ./.env.tmpl ./.env

.PHONY: infra
infra:
	@docker build \
		-f ./third_party/cartesi-coprocessor-operator/Dockerfile \
		-t operator:latest ./third_party/cartesi-coprocessor-operator
	@cd third_party/cartesi-coprocessor; docker compose -f docker-compose-devnet.yaml up --build

.PHONY: anvil
anvil:
	@cd third_party/cartesi-coprocessor; docker compose -f docker-compose-devnet.yaml up anvil --build

.PHONY: clean
clean:
	@rm -rf ./contracts/out
	@rm -rf ./contracts/cache
	@rm -rf ./contracts/broadcast
	@rm -rf ./third_party/cartesi-coprocessor/operator1-data
	@rm -rf ./third_party/cartesi-coprocessor/env/eigenlayer/anvil/devnet-operators-ready.flag
	@rm -rf ./output.car
	@rm -rf ./output.car.json
	@rm -rf ./output.cid
	@rm -rf ./output.size

.PHONY: test
test:
	@go test -p=1 ./... -coverprofile=./coverage.md -v
	@forge test --root contracts

.PHONY: coverage
coverage: test
	@go tool cover -html=./coverage.md

.PHONY: fmt
fmt:
	@go fmt ./...
	@forge fmt --root contracts

.PHONY: gen
gen:
	@go generate ./...

.PHONY: slot
slot:
	@rm -rf storage-layout
	@forge inspect ./src/SwapXHook.sol:SwapXHook storage-layout --root contracts >> storage-layout

.PHONY: v4
v4:
	@forge script ./contracts/script/V4Deployer.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		-vvvv

.PHONY: hook
hook:
	@forge script ./contracts/script/HookDeployer.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		-vvvv

.PHONY: demo
demo:
	@forge script ./contracts/script/Demo.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		-vvvv