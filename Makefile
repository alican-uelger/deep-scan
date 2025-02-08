.DEFAULT_GOAL := build

COVER_PROFILE := coverage.out
COVER_OUTPUT := coverage.html

.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -o ./dist/local/deep-scan ./main.go

.PHONY: check
check: check-int check-ext

.PHONY: check-int
check-int:
	-go fmt ./...
	-go vet ./...
	-govulncheck ./...

.PHONY: check-ext
check-ext:
	-golangci-lint run
	-goreleaser check

.PHONY: clean
clean:
	go clean
	-rm -rf ./dist
	-rm coverage.out
	-rm coverage.html

.PHONY: mocks
mocks:
	go run github.com/vektra/mockery/v2@v2.51.0

# build tags end2end, integration and unit. E.g.:
# //go:build end2end
# //go:build integration
# //go:build unit
#
# IntelliJ IDEA and GoLand integration:
# https://www.jetbrains.com/help/go/go-build.html
#
# Visual Studio Code integration:
# https://github.com/golang/vscode-go/wiki/settings#gobuildtags
.PHONY: test-e2e
test-e2e:
	go test \
		-cover \
		-coverpkg=./... \
		-coverprofile=${COVER_PROFILE} \
		-parallel=16 \
		-shuffle=on \
		-tags=end2end -v ./...
	go tool cover -html=${COVER_PROFILE} -o ${COVER_OUTPUT}

.PHONY: test-int
test-int:
	go test \
		-cover \
		-coverpkg=./... \
		-coverprofile=${COVER_PROFILE} \
		-parallel=16 \
		-shuffle=on \
		-tags=integration -v ./...
	go tool cover -html=${COVER_PROFILE} -o ${COVER_OUTPUT}

.PHONY: test-unit
test-unit:
	go test \
		-cover \
		-coverpkg=./... \
		-coverprofile=${COVER_PROFILE} \
		-parallel=16 \
		-shuffle=on \
		-tags=unit -v ./...
	go tool cover -html=${COVER_PROFILE} -o ${COVER_OUTPUT}

.PHONY: benchmark
# Run benchmarks and output to a file
benchmark:
	@echo "Running benchmarks..."
	@go test -tags="unit int end2end" -bench . -benchmem ./... | tee benchmark_results.txt
