BUILD := build/lambda-go

$(BUILD): main.go Makefile
	GOARCH=amd64 GOOS=linux go build -o $@
