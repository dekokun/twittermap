DEP := .bin/dep
SAMLOCAL := .bin/aws-sam-local

.PHONY: build
build:
	@$(MAKE) build/tweetget
	@$(MAKE) build/s3upload

build/tweetget: tweetget.go Makefile config.toml
	GOARCH=amd64 GOOS=linux go build -o $@ $<
	cp config.toml build/config.toml
build/s3upload: s3upload.go Makefile config.toml
	GOARCH=amd64 GOOS=linux go build -o $@ $<
	cp config.toml build/config.toml

.PHONY: install
install: $(DEP)
	$(DEP) ensure

.bin/%:
	@$(MAKE) setup-go
	@touch $@

.PHONY: setup-go
setup-go:
	GOBIN=$(abspath .bin) go get -v \
		github.com/golang/dep/cmd/dep \
