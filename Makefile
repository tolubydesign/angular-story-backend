STACK_NAME ?= serverless-go-demo
FUNCTIONS := call
# FUNCTIONS := get-products get-product put-product delete-product products-stream call
REGION := us-east-2

# To try different version of Go
GO := go

# Make sure to install aarch64 GCC compilers if you want to compile with GCC.
CC := aarch64-linux-gnu-gcc
GCCGO := aarch64-linux-gnu-gccgo-10

# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go

ci: build tests-unit

build:
	${MAKE} ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-${function})

build-%:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ${GO} build -o main ./functions/$*
	cp ./main $(ARTIFACTS_DIR)/.
	rm ./main
# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
# cd functions/$* && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 ${GO} build -o bootstrap ./
# go build -o api ./functions/call
# cd functions/$* && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 ${GO} build -o bootstrap

build-gcc:
	${MAKE} ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-gcc-${function})

build-gcc-%:
	cd functions/$* && GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=${CC} ${GO} build -o main

build-gcc-optimized:
	${MAKE} ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-gcc-optimized-${function})

build-gcc-optimized-%:
	cd functions/$* && GOOS=linux GOARCH=arm64 GCCGO=${GCCGO} ${GO} build -compiler gccgo -gccgoflags '-static -Ofast -march=armv8.2-a+fp16+rcpc+dotprod+crypto -mtune=neoverse-n1 -moutline-atomics' -o bootstrap

invoke:
# 	@sam local invoke --env-vars env-vars.json GetProductsFunction

# invoke-put:
# 	@sam local invoke --env-vars env-vars.json --event functions/put-product/event.json PutProductFunction

# invoke-get:
# 	@sam local invoke --env-vars env-vars.json --event functions/get-product/event.json GetProductFunction

# invoke-delete:
# 	@sam local invoke --env-vars env-vars.json --event functions/delete-product/event.json DeleteProductFunction

# invoke-stream:
# 	@sam local invoke --env-vars env-vars.json --event functions/products-stream/event.json DDBStreamsFunction

clean:
	@rm $(foreach function,${FUNCTIONS}, functions/${function}/bootstrap)

deploy:
# sam deploy --guided (-g means --guided)
	if [ -f samconfig.toml ]; \
		then sam deploy --stack-name ${STACK_NAME} --region ${REGION}; \
		else sam deploy -g --stack-name ${STACK_NAME} --region ${REGION}; \
  fi

tests-unit:
	@go test -v -tags=unit -bench=. -benchmem -cover ./...

tests-integ:
	API_URL=$$(aws cloudformation describe-stacks --stack-name $(STACK_NAME) \
	  --region $(REGION) \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text) go test -v -tags=integration ./...

tests-load:
	API_URL=$$(aws cloudformation describe-stacks --stack-name $(STACK_NAME) \
	  --region $(REGION) \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text) artillery run load-testing/load-test.yml

export GOBIN ?= $(shell pwd)/bin

STATICCHECK = $(GOBIN)/staticcheck

# Many Go tools take file globs or directories as arguments instead of packages
GO_FILES := $(shell \
	       find . '(' -path '*/.*' -o -path './vendor' ')' -prune \
	       -o -name '*.go' -print | cut -b3-)
MODULE_DIRS = .

.PHONY: lint
lint: $(STATICCHECK)
	@rm -rf lint.log
	@echo "Checking formatting..."
	@gofmt -d -s $(GO_FILES) 2>&1 | tee lint.log
	@echo "Checking vet..."
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && go vet ./... 2>&1) &&) true | tee -a lint.log
	@echo "Checking staticcheck..."
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && $(STATICCHECK) ./... 2>&1) &&) true | tee -a lint.log
	@echo "Checking for unresolved FIXMEs..."
	@git grep -i fixme | grep -v -e Makefile | tee -a lint.log
	@[ ! -s lint.log ]
	@rm lint.log
	@echo "Checking 'go mod tidy'..."
	@make tidy
	@if ! git diff --quiet; then \
		echo "'go diff tidy' resulted in chnges or working tree is dirty:"; \
		git --no-pager diff; \
	fi

$(STATICCHECK):
	cd tools && go install honnef.co/go/tools/cmd/staticcheck

.PHONY: tidy
tidy:
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && go mod tidy) &&) true