BUILD := build/lambda-go
DEP := .bin/dep
SAMLOCAL := .bin/aws-sam-local

$(BUILD): main.go Makefile config.toml
	GOARCH=amd64 GOOS=linux go build -o $@
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
