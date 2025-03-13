-include .env

.PHONY: env
env:
	@cp ./.env.tmpl ./.env

.PHONY: infra
infra:
	@cd third_party/cartesi-coprocessor; docker compose -f docker-compose-devnet.yaml --profile explorer up --build

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
	@cd contracts; forge script ./script/V4Deployer.s.sol --broadcast \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--verify \
		--verifier sourcify \
		--verifier-url http://localhost:5555 \
		-vvvv

.PHONY: hook
hook:
	@cd contracts; forge script ./script/HookDeployer.s.sol --broadcast \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--verify \
		--verifier sourcify \
		--verifier-url http://localhost:5555 \
		-vvvv

.PHONY: demo01
demo01:
	@forge script ./contracts/script/demo/01_demo.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--slow \
		-vvvv

.PHONY: demo02
demo02:
	@forge script ./contracts/script/demo/02_demo.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--slow \
		-vvvv

.PHONY: demo03
demo03:
	@forge script ./contracts/script/demo/03_demo.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--slow \
		-vvvv

.PHONY: demo04
demo04:
	@forge script ./contracts/script/demo/04_demo.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--slow \
		-vvvv

.PHONY: demo05
demo05:
	@forge script ./contracts/script/demo/05_demo.s.sol --broadcast \
		--root contracts \
		--rpc-url $(RPC_URL) \
		--private-key $(PRIVATE_KEY) \
		--slow \
		-vvvv

demo: demo01
	@sleep 2
	@$(MAKE) demo02
	@sleep 2
	@$(MAKE) demo03
	@sleep 2
	@$(MAKE) demo04
	@sleep 2
	@$(MAKE) demo05